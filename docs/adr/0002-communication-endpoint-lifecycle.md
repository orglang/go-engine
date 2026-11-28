# ADR-0002: Жизненный цикл Communication и его endpoints

- **Статус:** Accepted
- **Дата:** 2026-11-08

## Контекст

`Execution` является корнем агрегата. `Communication` — вложенная сущность, связывающая провайдерскую и клиентскую Execution для обмена сообщениями. Endpoint-ы представлены переменными Execution: liability-переменная провайдера и соответствующая asset-переменная клиента содержат один `CommID`. В рассматриваемом сценарии провайдерская Execution создаётся раньше клиентской.

Если откладывать создание Communication до появления клиента, можно заранее зарезервировать его идентификатор в liability-переменной провайдера. Это позволяет сохранить идентичность будущего соединения без преждевременного создания самой сущности.

## Рассмотренные варианты

### Вариант C: резервировать CommID в liability-переменной, создавать Communication вместе с client Execution

При создании provider Execution генерируется `CommID` и записывается в её liability-переменную, но запись Communication ещё не создаётся. При создании client Execution используется тот же `CommID`: в одной транзакции создаётся Communication с обоими endpoints и создаётся/привязывается liability-переменная клиента.

Преимущества:
- Provider Execution сразу получает стабильный `CommID`; её не нужно обновлять при создании клиента.
- Communication создаётся только тогда, когда известны оба endpoint, поэтому не требуется nullable `ClientRootRef`.
- Не нужно создавать Communication при первом Act или вводить отдельную операцию позднего обновления endpoint.

Недостатки и обязательные условия:
- До создания клиента `CommID` существует только как зарезервированное значение в liability-переменной; строки Communication ещё нет. Это должно быть допустимым состоянием модели и хранилища.
- Все операции, которым нужна именно сущность Communication (например, чтение состояния канала или запись Act), должны требовать существующую запись Communication и не трактовать один `CommID` как доказательство её существования.
- Создание Communication при создании клиента должно проверять, что `CommID` действительно зарезервирован provider Execution и ещё не использован; создание Communication и клиентской стороны должно быть атомарным.
- Нужно определить поведение, если client Execution так и не будет создана: зарезервированный `CommID` останется в provider Execution без соответствующей записи Communication.

### Вариант A: создавать Communication вместе с provider Execution

При создании provider Execution в той же транзакции создаются Communication и liability-переменная провайдера. Communication сразу получает `ProviderRootRef` и свой `CommID`; `ClientRootRef` остаётся пустым до создания client Execution. При создании клиента существующий Communication дополняется второй ссылкой.

Преимущества:
- `CommID` доступен провайдеру с момента создания его Execution и может использоваться всеми последующими операциями провайдера.
- Liability-переменная провайдера сразу имеет окончательную связь с Communication; не требуется позднее обновлять уже созданную переменную.
- Создание клиента сводится к присоединению второго endpoint к уже известной сущности.

Недостатки:
- Между созданием provider и client Execution существует Communication с незаполненным `ClientRootRef`.
- Хранилище и прикладные операции должны допускать это состояние, но не разрешать Act до появления клиента.

### Вариант B: создавать Communication вместе с client Execution

При создании provider Execution Communication ещё не существует, а её liability-переменная не имеет `CommID`. В транзакции создания client Execution создаётся Communication с обоими endpoints, после чего обновляется liability-переменная уже существующей provider Execution и создаётся/привязывается переменная клиента.

Преимущества:
- Communication появляется только тогда, когда известны обе Execution; не требуется состояние с незаполненным `ClientRootRef`.
- Создаваемый Communication сразу имеет оба endpoint.

Недостатки:
- До создания клиента у провайдера нет `CommID`; любые операции, которым нужна идентичность соединения, должны быть отложены или иметь отдельное представление состояния «ещё не подключён».
- Создание клиента должно атомарно изменять уже существующую provider Execution и её liability-переменную вместе с созданием Communication и client Execution. Это расширяет транзакционную границу и требует найти и обновить провайдерскую сторону.
- Появляется поздняя мутация ранее созданной provider Execution ради завершения её связи.

## Решение

Запись Communication хранит обе endpoint-ссылки явно: `ProviderRootRef` и `ClientRootRef`. В `pool_exec_comms` и `proc_exec_comms` соответствующие `provider_exec_id`/`provider_exec_rn` и `client_exec_id`/`client_exec_rn` обязательны. Pool Communication и Process Communication — независимые сущности со своими записями и жизненными циклами. Process Communication не создаётся и не копируется из Pool Communication; её создание при присоединении второй Process Execution рассматривается отдельно в [go-workspace issue #112](https://github.com/orglang/go-workspace/issues/112). И provider, и client могут инициировать как PubRec, так и SubRec; тип события не определяет сторону-инициатора. Каждый Act добавляется в общий поток Communication с увеличением общего номера Act.

Для Pool Communication генерировать `CommID` при создании provider Execution и сохранять его в liability-переменной. Сам Pool Communication создавать в транзакции создания client Execution, когда известны оба endpoint. Это решение о Pool Communication не задаёт момент создания Process Communication: process-side attachment lifecycle вынесен в [go-workspace issue #112](https://github.com/orglang/go-workspace/issues/112).

### Владение Communication и конкурентные изменения

`Communication` является вложенной сущностью агрегата provider Execution. Provider Execution — единственный aggregate root, через который выполняются все изменения соответствующего Communication, включая добавление Act и изменение общего состояния взаимодействия. Это правило не зависит от того, какая сторона инициировала операцию: client может инициировать PubRec или SubRec, но запись Communication всё равно изменяется в транзакции с проверкой версии provider Execution.

Каждая операция, изменяющая Pool или Process Communication, должна проверять именно ту версию `provider.exec_rn`, которую операция прочитала вместе со snapshot Communication, и атомарно увеличивать `provider.exec_rn` в той же транзакции, что и изменение Communication. Нельзя перечитывать актуальную версию provider root только внутри write-транзакции: это скрыло бы конкурентное изменение, произошедшее после чтения snapshot. Если ожидаемая версия устарела, CAS завершается конфликтом и вся транзакция откатывается. Поэтому две разные Execution, конкурирующие за один Pool или Process Communication, конфликтуют на версии одного и того же provider aggregate root. Отдельная версия Communication для optimistic concurrency control не вводится. Для Pool и Process Communication `act_nr` хранит ревизию provider Execution, которая была получена CAS-операцией и породила Act (`observed provider.exec_rn + 1`); тем самым `act_nr` одновременно задаёт порядок Act и не требует отдельного инкремента. Сам concurrency control по-прежнему выполняется исключительно CAS на `provider.exec_rn`, а не сравнением `act_nr`. За одну операцию добавляется не более одного Act. `offset_nr` обновляется под защитой CAS provider root.


### Доступ к Communication со стороны клиента при партицировании

Если таблица Pool Communication будет партицирована по владельцу агрегата (`provider_exec_id`), запись Communication естественно размещается в партиции provider Execution. Это соответствует владению и позволяет выполнять изменения Communication рядом с provider aggregate root. Однако такой ключ партицирования создаёт отдельную задачу для чтения со стороны клиента: клиентская Execution знает `CommID` из своей asset-переменной, но одного `CommID` может быть недостаточно, чтобы определить партицию и эффективно выполнить partition pruning.

При этом глобальную уникальность `CommID` нельзя автоматически гарантировать обычным `UNIQUE(comm_id)` на PostgreSQL partitioned table, если ключ партицирования не входит в ограничение уникальности. Ограничение вида `UNIQUE(provider_exec_id, comm_id)` обеспечивает уникальность только в сочетании с ключом владельца и означает, что для адресации Communication клиенту потребуется также знать `provider_exec_id`. Это значение можно хранить в клиентской endpoint-переменной/ссылке вместе с `CommID`.

Если контракт должен позволять находить Communication только по глобальному `CommID`, потребуется отдельный глобально адресуемый lookup-механизм, например небольшая непартицированная registry-таблица `CommID → provider_exec_id`. Такой registry добавляет ещё одну сущность и требует атомарно поддерживать её при создании и удалении Communication. Альтернативой может быть отдельная непартицированная таблица идентификаторов, если её свойства и ограничения подходят модели.

Таким образом, партицирование по `provider_exec_id` не меняет aggregate ownership или правило OCC, но влияет на контракт адресации Communication с клиентской стороны. До выбора физической схемы необходимо явно решить, будет ли клиентская сторона передавать `provider_exec_id` вместе с `CommID` или система обязана поддерживать lookup только по `CommID`. Это архитектурное решение о способе поиска; оно не должно приводить к переносу Communication в клиентскую партицию или к изменению его владельца.

В рамках PR #61 тот же принцип aggregate ownership применяется к Process Communication независимо от её отдельного lifecycle: provider Execution является единственным aggregate root, и любое изменение Process Communication сериализуется CAS по ожидаемой `provider.exec_rn`. Это не связывает жизненные циклы Pool и Process Communication и не предполагает копирования их состояния. Attachment Process Communication остаётся отдельной задачей (см. [go-workspace issue #112](https://github.com/orglang/go-workspace/issues/112)).

`client.exec_rn` увеличивается только тогда, когда операция одновременно изменяет собственное состояние client Execution (например, её переменные). Сам факт изменения общего Communication не является основанием увеличивать ревизию client Execution.

При создании клиентской Execution система получает зарезервированный `CommID` из provider liability-переменной и в той же транзакции создаёт запись Communication и asset-переменные клиента. Liability-переменная провайдера и asset-переменная клиента тем самым образуют два endpoint-а одной Communication и ссылаются на один `CommID`. До этой транзакции строки Communication нет.

Оба endpoint должны быть определены до первого Act. Первый Act не создаёт Communication и не выполняет отложенную привязку endpoints.

## Повторная оценка выбора

Вариант C следует предпочесть варианту A, если модель допускает зарезервированный `CommID` без строки Communication до создания клиента. Он сохраняет раннюю и стабильную идентичность для provider Execution, но создаёт сам Communication только после появления обоих endpoints. Это устраняет как позднее изменение provider Execution из варианта B, так и неполный Communication из варианта A.

Если же по инвариантам хранения каждый `CommID` в liability-переменной обязан немедленно обозначать существующую запись Communication, вариант C неприменим; тогда остаётся вариант A.

## Инварианты

- `CommID` глобально идентифицирует Communication; отдельный составной reference для него не вводится.
- Communication явно хранит `ProviderRootRef` и `ClientRootRef`; ссылки включают идентификатор и ревизию Execution.
- Pool Communication принадлежит агрегату provider Execution; provider Execution является единственным aggregate root для его изменений, а Pool Communication не является самостоятельным корнем агрегата.
- Любое изменение Pool или Process Communication проверяет ожидаемый `provider.exec_rn` из прочитанного snapshot и атомарно увеличивает `provider.exec_rn`; отдельная версия Communication и CAS по `offset_nr` для optimistic concurrency control не используются.
- `client.exec_rn` увеличивается только при изменении собственного состояния client Execution, а не только из-за изменения Pool Communication.
- Для Pool и Process Communication `act_nr` равен новой ревизии provider Execution, полученной CAS-операцией, которая добавляет Act. Он задаёт порядок Act, но не является самостоятельным concurrency token; отдельный инкремент счётчика не выполняется. Одна операция добавляет не более одного Act.
- При выбранном варианте C запись Communication создаётся только после появления обоих endpoints; принадлежность endpoint-ов определяется переменными Execution, содержащими общий `CommID`.
- До создания клиента provider liability-переменная может содержать зарезервированный `CommID`, для которого запись Communication ещё отсутствует.
- Для одной пары endpoints используется один Communication и один `CommID`.
- Pool Communication и Process Communication не создаются друг из друга и не зависят друг от друга. Момент создания Process Communication не входит в объём этого PR; он должен быть определён отдельной process-side attachment операцией (см. [go-workspace issue #112](https://github.com/orglang/go-workspace/issues/112)). После создания любое изменение Process Communication, включая offset и Acts, защищается CAS по `provider.exec_rn`. `offset_nr` — только позиция в потоке, не concurrency token. PubRec и SubRec от любой из сторон используют одну process-side запись Communication и общий счётчик Acts; тип события не ограничивает сторону-инициатора.
- Создание Communication с обоими endpoints и клиентской стороны выполняется атомарно с созданием клиентской Execution и её переменных.
- Первый Act допустим только после заполнения обоих endpoints.
- `ChnlID` остаётся отдельной идентичностью interaction channel и не заменяет `CommID`.

## Последствия

### Положительные

- Провайдер получает стабильный `CommID` сразу при создании своей Execution.
- Обе стороны используют одну запись Communication и одну идентичность канала.
- Communication создаётся сразу с обоими endpoint; не требуется отложенная привязка.
- Жизненный цикл Communication согласован с тем, что provider Execution является его единственным aggregate root.
- Конкурирующие изменения одного Communication сериализуются optimistic concurrency control на `provider.exec_rn`, даже если операции инициированы разными endpoint-ами.

### Отрицательные

- До создания клиента существует зарезервированный `CommID` без строки Communication.
- Операции, которым нужна сущность Communication, должны отличать зарезервированный идентификатор от созданного Communication.
- Все операции изменения Communication, включая инициированные client-стороной, должны иметь доступ к provider aggregate root и выполнять запись в его транзакционном контексте.

## Когда пересмотреть решение

Пересмотреть ADR следует, если изменится порядок создания сторон, появится сценарий, в котором Communication должен существовать без провайдерской Execution, либо будет введён отдельный жизненный цикл установления соединения, независимый от создания Execution.
