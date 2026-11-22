# ADR-0004: Модель корней агрегатов и semantic references

- **Статус:** Proposed (целевой Description aggregate и reference contracts; открытые предметные решения перечислены в §23)
- **Дата:** 2026-11-22
- **Связано:** go-workspace#128, go-engine#188, go-workspace#122

## Цель

Зафиксировать фактическую модель агрегатов в `proc` и `pool`, отделить её от ADT и подготовить обоснованные решения о reference-контрактах и persistence roots. Этот документ задаёт целевую предметную модель. Публичные контракты и схема БД могут быть изменены для её реализации; конкретные изменения определяются после согласования §23 и migration plan. Предварительные варианты и вопросы в §7, §12–15 являются историей исследования; целевое предложение и открытые решения находятся в §19–24.

## 1. Оси модели

- **Контекст исполнения:** `proc` или `pool`. Контексты имеют отдельные application/core пакеты, Ent clients и физические таблицы.
- **Вид агрегата:** type, term, computation. Вид определяет бизнес-идентичность и специализированные инварианты.
- **ADT:** набор семантических value types, primitive identity/revision types и semantic ports. ADT не является ещё одним контекстом исполнения и не должен владеть конкретными persistence adapters.

## 2. Предварительная матрица фактически реализованных агрегатов

| Контекст | Type | Term | Computation |
|---|---|---|---|
| `pool` | `pool/typedef` → `pool_type_defs` | `pool/termdef` → `pool_term_decs` | `pool/compexec` → `pool_comp_execs` |
| `proc` | `proc/typedef` → `proc_type_defs` | `proc/termdec` → `proc_term_decs` | `proc/compexec` → `proc_comp_execs` |

`proc/termdef` присутствует в исходниках, но текущий core содержит незавершённые заглушки; в initial SQL migration и Ent schema соответствующего process term-definition root нет. Это описание current state, а не целевой модели: proc.termdef входит в целевой Term Description и должен быть реализован как часть этого направления.

Таблицы и имена в матрице описывают текущую реализацию, а не целевую схему.

## 3. Текущие reference-типы

| ADT package | Go type | Поля | Использование |
|---|---|---|---|
| `adt/typesem` | `SemRef` | `TypeID identity.ADT`, `TypeRN seqnum.ADT` | Pool и Process Type |
| `adt/termsem` | `SemRef` | `TermID identity.ADT`, `TermRN seqnum.ADT` | Pool Term и Process Term |
| `adt/compsem` | `SemRef` | `CompID identity.ADT`, `CompRN seqnum.ADT` | Pool и Process Computation |

Наблюдение: текущие типы reference уже разделены по виду агрегата, а не по декартовой комбинации контекста и вида. При этом `proc` и `pool` используют одинаковый Go reference type для одного вида, но отдельные ID и физические таблицы.

Предварительная гипотеза: контекст исполнения не должен входить в reference type, если reference идентифицирует экземпляр агрегата независимо от места исполнения. Проверить это на API/wire контрактах, переносимости ссылок и инвариантах.

## 4. Различия контрактов и поведения

- `typesem.Repo` предоставляет `TouchRef` и `GetRefsByQNs`; type reference участвует в lookup по QN и optimistic revision flow в Process Type.
- `compsem.Repo` предоставляет `TouchRef`; computation revision участвует в исполнении и синхронизации состояний.
- `termsem` сейчас содержит value type и конструктор, но не аналогичный semantic repository port.
- `proc/typedef` реализует Create/Modify/Retrieve и optimistic revision check.
- `pool/typedef` в текущем core предоставляет Create.
- `proc/termdec` поддерживает Incept/Create/RetrieveSnap/RetreiveRefs; `pool/termdef` в текущем core предоставляет Create.
- `proc/compexec` и `pool/compexec` имеют различные use cases и execution semantics.

Следовательно, одинаковые поля reference не являются достаточным основанием для объединения всех semantic packages или repository ports. Решение по `typesem` и `termsem` должно опираться на инварианты и реальные потребители.

## 5. Предварительная карта persistence roots и зависимых данных

- Type roots: `proc_type_defs`, `pool_type_defs`; Ent schema задаёт `type_id`, `type_rn`, `exp_ck`.
- Term roots: `proc_term_decs`, `pool_term_decs`; Ent schema задаёт `term_id`, `term_rn`, `liab_var`, `asset_vars`.
- Computation roots: `proc_comp_execs`, `pool_comp_execs`; Ent schema задаёт `comp_id`, `comp_rn`, `liab_mode`.
- Computation variable tables содержат `comp_id` и `comp_rn`; PostgreSQL inheritance используется для struct/linear variable tables.
- Connection tables содержат явные provider/client computation references (ID + revision); connection action tables содержат computation ID + revision.
- Type expressions хранятся отдельно и адресуются через `exp_ck`; их принадлежность и lifecycle нужно классифицировать как часть aggregate state либо как отдельную shared immutable сущность.
- Descriptor binding tables (`proc_desc_binds`, `pool_desc_binds`) связывают QN с ID и kind; необходимо отдельно решить, являются ли они частью aggregate consistency boundary или отдельным индексом/registry.

В initial migration нет отдельных root tables поверх type/term/computation tables: эти таблицы сами содержат ID и revision. Это предварительно поддерживает модель «specialized root table is the aggregate root», но не доказывает, что все lifecycle и referential-integrity требования удовлетворены.

## 6. Известные несоответствия, требующие проверки

- `typesem/core.go` содержит `typeKind` и `PoolKind`/`ProcKind`, которые не отражены в текущем reference contract; проверить все usages перед удалением.
- `proc/termdef` выглядит как незавершённый/legacy пакет и использует `typesem.SemRef`, тогда как фактический process term aggregate находится в `proc/termdec` и использует `termsem.SemRef`.
- `proc/termdec/adapter/ds_ent_pgx.go` содержит запросы к `dec_roots` и `desc_sems`; эти таблицы не объявлены в единственной обнаруженной migration `db/postgres/owner/000001_initial.up.sql`. Установить, являются ли запросы устаревшими, схема неполной или объекты управляются вне этого migration tree. Не менять DDL до выяснения.
- В `proc/typedef/adapter/ds_pgx.go` обнаружены SQL-фрагменты со старыми/нестыкующимися именами колонок; проверить, используются ли эти методы текущим DI wiring и актуальны ли запросы.

## 7. Решения, ожидающие завершения аудита

1. Терминология: оставить `SemRef`, перейти к `TypeRef`/`TermRef`/`CompRef` или использовать `AggRef` с типизированными обёртками.
2. Состав reference: ID и revision; роль QN (вероятно, не часть идентичности); инварианты revision для каждого вида.
3. Root model: оставить существующие специализированные root tables или ввести отдельную root table только при доказанной необходимости для целостности/lifecycle.
4. Таблицы aggregate state: для каждой зависимой таблицы подтвердить прямую ссылку на root Ref или обосновать исключение.
5. `typesem` и `termsem`: объединять только если общими являются не только поля, но и бизнес-инварианты/операции.
6. ADT ports: определить, какие semantic ports действительно принадлежат ADT, и не протаскивать Ent/persistence DTO в core.

## 8. Ограничения реализации

- Не менять SDK в рамках реализации issue #128; это ограничение scope, а не требование сохранять прежний wire/API contract в go-engine.
- Не менять SQL migrations без отдельного обоснования и совместимого migration plan.
- Не создавать все декартовы комбинации автоматически.
- Не вводить универсальный CRUD или универсальный reference только ради уменьшения числа типов.
- Не переносить Ent API или persistence DTO за границы adapters.
- Не сохранять внешние wire/API contracts только ради обратной совместимости; API проектируется под целевые domain/reference contracts.


## 9. Проверка по публичной TypeSpec-модели

В текущем `tsp` уже существует общий namespace `OrgLang.Any` с тремя типизированными ссылками:

- `Any.TypeRef { type_id, type_rn }`;
- `Any.TermRef { term_id, term_rn }`;
- `Any.CompRef { comp_id, comp_rn }`.

Операции создания Process и Pool Type возвращают один `Any.TypeRef`; операции создания Process/Pool Declaration возвращают один `Any.TermRef`; списки Process/Pool executions возвращают `Any.CompRef[]`. Следовательно, публичный wire-контракт уже подтверждает разделение ссылок по виду агрегата и не кодирует execution context в типе ссылки.

Это не означает, что сами агрегаты или их application services общие. TypeSpec сохраняет контекстные модели (`Proc.Type.DefSnap`, `Pool.Type.DefSnap`, `Proc.CompSnap`, `Pool.CompSnap`), а физические roots остаются раздельными.

## 10. Уточнение по Computation и Connection

ADR-0002 фиксирует, что Pool Connection и Process Connection — независимые вложенные сущности со своими жизненными циклами. Обе endpoint-ссылки представлены явно как provider/client computation refs. При этом Connection принадлежит provider Computation aggregate: изменения Connection выполняются через provider root и защищаются CAS по его revision, независимо от того, какая сторона инициировала действие.

Вывод: общий `CompRef` обозначает computation identity/revision и может использоваться в обоих контекстах, но сам по себе не выражает владение, право записи, транзакционную границу или lifecycle Connection. Эти гарантии должны оставаться в контекстных core/use-case контрактах.

## 11. Уточнение по TypeExp

Хотя обе стороны используют одинаковый ключ `exp_ck` и отдельные таблицы `proc_type_exps` / `pool_type_exps`, их payload schemas уже различаются: Process содержит, например, `plus[]`/`with[]`, Pool — `plus`/`with` и дополнительные `up`/`down`. Поэтому общий `TypeRef` не является основанием объединять TypeExp state или его persistence model. На текущем уровне это context-owned state, адресуемый из соответствующего Type root через `exp_ck`.

## 12. Рабочее решение для следующего этапа

1. Current state имеет три семейства ссылок по виду агрегата: Type, Term и Computation. Target contracts определены в §20; старые SemRef не являются ограничением.
2. Не включать `Proc`/`Pool` в тип ссылки. Контекст остаётся частью владельца, репозитория, application service и persistence boundary.
3. Не объединять Proc/Pool roots, TypeExp state, Connection lifecycle или use cases только из-за общей формы ссылки.
4. Не вводить универсальный `AggRef`: публичная модель уже имеет три конкретных типа, и они точнее выражают допустимые связи.
5. Перед изменением Go contracts проверить все conversion sites, сериализацию, DI wiring и потребителей semantic ports; отдельно определить, должна ли revision быть частью каждой ссылки в каждом use case.


## 13. Отдельная гипотеза: схлопнуть TypeRef и TermRef

Эта гипотеза рассматривается отдельно от вопроса об объединении aggregate roots. У `typesem.SemRef` и `termsem.SemRef` одинаковая физическая форма: пара `identity.ADT` + `seqnum.ADT`, одинаковый конструктор `New()` и одинаковая базовая трактовка revision. Это аргумент в пользу общего value representation.

При этом обнаружены реальные семантические различия:

- Type references разрешаются по QN через `typesem.Repo.GetRefsByQNs`; этот port также предоставляет `TouchRef`. `proc/termdec` использует его, чтобы разрешить типы liability/asset variables.
- Term references сейчас не имеют аналогичного semantic repository port. Term QN разрешается контекстными `termdef`/`termdec` repositories и descriptor bindings.
- `termvar.VarRec.TypeRef` всегда ссылается на Type, тогда как `pool/termexp.SpawnSpec2.ProcTermRef` ссылается на Process Term. В этих местах вид ссылки — часть предметного смысла поля.
- Descriptor binding (`descsem`) различает `TypeKind` и `TermKind`, то есть реестр QN уже хранит явный discriminator вида сущности.
- Публичный TypeSpec API в настоящее время экспортирует `Any.TypeRef` и `Any.TermRef` отдельно. Это факт current state, а не ограничение: целевой API может быть переработан под новые typed references.

### Варианты модели

1. **Общий внутренний value type, разные именованные поля/контракты.** Ввести общий `adt/semref.Ref` (ID + revision), а Type/Term boundary packages использовать для семантических операций и явных конвертаций. Это устраняет структурное дублирование, но требует не допускать неявной подмены Type ссылкой на Term.
2. **Один полностью общий `SemRef` без discriminator.** Самый компактный вариант, но вызывающий код должен знать вид root из контекста. Он слабее защищает от ошибочного lookup в неправильном repository и не подходит для автономной передачи/хранения ссылки без контекста.
3. **Общий `SemRef` с kind discriminator.** Ссылка становится самодостаточной для маршрутизации, но kind дублирует часть информации, уже заданной типизированным полем/контрактом и descriptor registry. Также нужно определить, входит ли kind в идентичность или только в маршрутизацию.
4. **Оставить TypeRef и TermRef независимыми.** Сохраняет максимальную статическую типобезопасность, но дублирует одинаковую пару ID/revision и конвертеры.

### Предварительная рекомендация

Ранний вариант из инвентаризационной части не является решением. Целевая модель из §19–20 задаёт единый Description identity space внутри контекста при сохранении разных Type/Term domain contracts. Старые TypeSpec `Any.TypeRef`/`Any.TermRef` не являются обязательными wire contracts.

До реализации проверить: (а) уникальность identity между Type и Term внутри одного persistence context; (б) все сериализационные и generated conversion paths; (в) правила revision/CAS для Term roots и Type roots; (г) возможность случайно передать общий Ref в repository неправильного вида; (д) не требует ли автономная ссылка kind для корректной маршрутизации. Если хотя бы один из этих контрактов требует статической гарантии вида, оставить TypeRef/TermRef как distinct named wrappers поверх общей пары.


### Результат проверки revision и маршрутизации

Предварительное сходство Type/Term по паре полей не означает одинаковые revision guarantees:

- Process Type реализует optimistic modification: `Modify` сравнивает входящий `TypeRN` с сохранённым и увеличивает revision; `typesem.Repo.TouchRef` выполняет CAS update по ID + revision.
- Process Term на текущем пути создаётся с `TermRN`, но `RetrieveSnap` в `proc/termdec/adapter/ds_ent_pgx.go` выбирает запись только по `TermID` и не проверяет `TermRN`. В просмотренных core/API нет операции изменения Term или CAS-порта, аналогичного `typesem.Repo.TouchRef`.
- Pool Type и Pool Term в текущем core имеют только `Create`; для них не обнаружено симметричного modify/CAS use case.

Значит, revision сейчас присутствует в Type/Term references как часть wire/data shape, но его фактические гарантии неодинаковы и зависят от контекста/use case. До унификации нужно решить, является ли revision обязательной частью общего Ref для всех видов корней или только version token там, где поддерживается optimistic concurrency. В частности, нельзя заявлять, что общий тип автоматически даёт одинаковую семантику CAS.

`identity.New()` использует xid, поэтому коллизии между разными таблицами практически маловероятны, но схема не задаёт единого cross-table uniqueness constraint. Общий Ref без kind discriminator сам по себе не позволяет маршрутизировать ID в Type или Term repository; ожидаемый вид должен оставаться известен из типизированного поля/use case либо кодироваться явно.

Это усиливает рекомендацию: общий низкоуровневый `Ref{ID, Revision}` возможен, но семантический контракт Type и Term пока следует оставить раздельным. Полное схлопывание до одного неразмеченного `SemRef` преждевременно, пока не согласованы revision policy и правила маршрутизации.


## 14. Пересмотр модели: единый aggregate для Type и Term

Этот раздел фиксирует уточнения, полученные при обсуждении issue #128. В случае противоречий он имеет приоритет над предварительными гипотезами разделов 2, 6, 12 и 13.

### 14.1. Контексты и сущности агрегата

\`proc\` и \`pool\` — бизнес-контексты (bounded contexts), а не разновидности одного reference. Каждый контекст владеет своими экземплярами агрегатов, application/use-case контрактами и persistence boundary.

Type и Term следует рассматривать не обязательно как два самостоятельных агрегата, а как **две разные сущности первого уровня одного семантического агрегата** (рабочее имя — \`descsem\`). Это не означает стирания различий между Type и Term: у сущностей остаются собственные kind, поля, правила интерпретации, операции и типизированные связи. Гипотеза объединения касается общей границы идентичности, ревизий, ссылочной целостности и жизненного цикла descriptor semantics.

Следовательно, матрица должна описывать размещение сущностей/ролей внутри агрегатов по контекстам, а не автоматически постулировать отдельный root aggregate для каждой ячейки. Отдельно требуется проверить, какие из Type/Term сущностей являются aggregate roots, а какие — элементы общего \`descsem\` root.

\`proc/termdef\` — развиваемая часть целевой модели, а не legacy-контракт. Отсутствие завершённого core API или отдельной persistence-модели в текущем состоянии означает незавершённую реализацию; оно не является аргументом против включения Process Term Definition в целевую архитектуру.

### 14.2. Общая модель ссылки и RN

Рабочая гипотеза: Type и Term имеют общую семантику идентичности, revision и ссылочной целостности. Поэтому схлопывание следует исследовать как объединение Type и Term в один aggregate (условно \`descsem\`), внутри которого Type и Term остаются двумя различными сущностями первого уровня, а не только как переиспользование общей структуры Go.

В терминологии модели использовать \`RN\` — revision number (\`revnum\`), а не \`Revision\`.

### 14.3. Zero-value RN как selector чтения

Вместо отдельного `RetrieveLatest(ID)` можно рассмотреть один контракт чтения с обычным value type для RN:

- `ID != zero`, `RN != zero`: запросить точную ревизию; отсутствие совпадающего текущего состояния должно давать явный результат/ошибку.
- `ID != zero`, `RN == zero`: запросить актуальную ревизию сущности.
- `ID == zero`: некорректный запрос.

В Go RN не обязательно делать nillable: отсутствие явно заданной ревизии кодируется zero-value типа `revnum.ADT`. Это позволяет сохранить value semantics и не вводить указатели и nil-проверки в query/reference structures. Ненулевой RN задаёт exact-revision semantics.

Важно разделять **durable reference** и **read selector**. Если durable reference должна обозначать конкретное состояние, её RN должен быть ненулевым и проверяться при создании/десериализации. Zero RN допустим как режим запроса «разрешить актуальную ревизию», но не как полноценная сохранённая ссылка на конкретное состояние. Один и тот же Go value type может использоваться в обоих контрактах только при явном различении этих инвариантов на уровне core/use case.

«Точная ревизия» здесь означает проверку, что RN совпадает с ревизией текущего сохранённого состояния. Это не обещание чтения исторического snapshot; для последнего потребуется отдельная история версий или versioned storage.

Контракт следует определить на уровне core/port, а не оставлять неявным поведением adapter. Тесты должны покрывать точное чтение, чтение latest при zero RN, отсутствующий ID и ненулевой RN без совпадающего состояния. Также нужно определить поведение при RN старше текущего и при удалённой/неизвестной ревизии.

### 14.4. Уточнение текущих несоответствий

- Отсутствующие revision-операции в \`termsem\` и непоследовательная проверка RN в \`proc/termdec\` трактуются как незавершённость реализации, а не как доказательство иной семантики Type и Term.
- Запросы к \`dec_roots\` и \`desc_sems\` в \`proc/termdec\` считаются устаревшими. Их не следует трактовать как свидетельство существования внешней/неучтённой схемы без новых данных.
- Текущий факт, что \`RetrieveSnap\` фильтрует только по \`TermID\`, остаётся наблюдением о реализации; целевой контракт должен быть определён явно через exact RN или nil RN semantics.

### 14.5. Следующие вопросы аудита

1. Определить границу и состав общего \`descsem\` aggregate: Type, Term, descriptor/QN bindings, их revision и lifecycle.
2. Определить, является ли один \`descsem\` root общим для всего \`proc\`/\`pool\` контекста или каждый контекст имеет собственный root и собственную identity space.
3. Сопоставить ссылки и foreign keys во всех Type/Term consumers с новой aggregate boundary; установить, какие ссылки должны нести RN, а какие могут запрашивать latest.
4. Определить durable reference contract отдельно от read selector contract и выбрать Go representation для optional RN.
5. Включить будущий \`proc/termdef\` в целевую матрицу, не выдавая текущие заглушки за реализованный persistence root.


## 15. Аудит существующего descsem (фактическая реализация)

Проверка исходников показывает, что internal/adt/descsem уже существует. Однако текущий descsem пока **не является полноценным aggregate root для Type и Term**. Это минимальный descriptor/QN binding registry:

- descsem.SemRec содержит только DescQN, DescID и Kind.
- Kind различает TypeKind и TermKind.
- descsem.Repo предоставляет только AddRec; операций retrieve, resolve, modify, delete или revision/CAS нет.
- PGX adapter выполняет только INSERT INTO <context table> (desc_qn, desc_id, kind).
- В миграции proc_desc_binds и pool_desc_binds имеют desc_qn ltree UNIQUE, desc_id varchar, kind smallint; revision в этих таблицах нет.
- Внешний ключ от desc_id к type/term storage и constraint на согласованность kind с целевой таблицей в initial migration не заданы.

Следовательно, название descsem уже занято общей семантической регистрацией descriptor names, но текущая модель пока не выражает предлагаемую общую identity/revision/lifecycle boundary. Это полезная отправная точка и, возможно, естественное место эволюции, но не доказательство, что агрегат уже реализован.

### 15.1. Текущее использование по контекстам

| Use case | Создаёт descriptor binding | Kind | Данные сущности |
|---|---|---|---|
| proc/typedef.Create | Да | TypeKind | proc_type_defs + proc_type_exps |
| pool/typedef.Create | Да | TypeKind | pool_type_defs + pool_type_exps |
| pool/termdef.Create | Да | TermKind | pool_term_decs |
| proc/termdec.Incept/Create | Нет в текущем core path | — | proc_term_decs |
| будущий proc/termdef | Целевое поведение ещё предстоит определить | Вероятно TermKind | Отдельная целевая модель не реализована |

У proc/termdec DI wiring также не получает descsem.Repo, а его текущий adapter содержит старые запросы к dec_roots/desc_sems. По принятому в обсуждении уточнению эти SQL-запросы считаются устаревшими; не следует выводить из них дополнительную схему. При этом отсутствие descriptor binding на текущем process term path — реальный пробел относительно Pool Term и Type paths, который нужно закрыть в целевой модели.

### 15.2. QN resolution сейчас распределён по потребителям

- Process Type lookup по QN идёт через proc_desc_binds с последующим чтением proc_type_defs.
- Pool Type lookup по QN идёт через pool_desc_binds с последующим чтением pool_type_defs.
- Pool Term lookup по QN идёт через pool_desc_binds с последующим чтением pool_term_decs.
- Process Term lookup в текущем proc/termdec adapter опирается на устаревшие SQL fragments и не является надёжным образцом целевого контракта.
- typesem.Repo.GetRefsByQNs — не общий descriptor resolver: он жёстко соединяет proc_desc_binds с proc_type_defs и возвращает только Type refs.

Иными словами, descsem пока пишет QN/kind binding, но чтение и проверка kind/сущности реализованы в контекстных специализированных repositories. Это создаёт дублирование lookup logic и не обеспечивает единую ссылочную целостность.

### 15.3. Архитектурная гипотеза для целевой модели

Исследовать эволюцию descsem в общий aggregate root внутри каждого bounded context (proc и pool), который владеет descriptor identity, QN, kind и RN. Type и Term в таком случае остаются разными сущностями/вариантами первого уровня одного агрегата; специализированные данные Type (включая TypeExp) и Term хранятся как kind-specific state, принадлежащий этому root.

Это пока гипотеза, а не принятое решение. Необходимо сравнить как минимум две формы:

1. **Единый Description root + специализированное состояние.** Один Description root на descriptor ID, содержащий QN, kind и RN; type/term-specific tables являются state/parts этого root и используют общий ID/RN.
2. **Общий descriptor registry над самостоятельными Type/Term roots.** descsem централизует QN и kind, но Type/Term сохраняют собственные root identities/revisions и consistency boundaries.

Первая форма соответствует сильному прочтению идеи «две сущности первого уровня одного агрегата» и может устранить дублирование identity/revision semantics. Вторая — более эволюционный вариант, но сама по себе не объединяет aggregate boundary и не гарантирует атомарность cross-kind invariants.

### 15.4. Существенные вопросы до выбора схемы

1. Решено в §17: отдельный Description aggregate внутри каждого bounded context; общей identity space между proc и pool нет.
2. Является ли QN частью identity root или изменяемым unique name/binding, который можно переименовать независимо от ID?
3. Должен ли RN жить в Description root и версионировать изменения Type/Term state совместно?
4. Как моделировать TypeExp: часть Type variant state под тем же root RN или отдельный immutable/versioned child resource?
5. Как обеспечить, чтобы kind соответствовал ровно одному специализированному состоянию и чтобы QN разрешался в согласованный ID/RN?
6. Должны ли Process Term Declaration (termdec) и Process Term Definition (termdef) быть двумя use cases одной Term-сущности/варианта, либо это разные роли/сущности внутри descriptor aggregate?
7. Как спроектировать миграцию текущих таблиц и API на целевые typed references без требования wire compatibility?

### 15.5. Наблюдение по zero RN

В internal/adt/seqnum уже определён seqnum.Zero, а seqnum.ADT является числовым value type. Это позволяет выразить selector «latest» без pointer/nil: zero RN означает resolve current; ненулевой RN означает проверку соответствия текущей ревизии. Это должно оставаться семантикой query/use case, а не означать, что durable reference может содержать нулевую ревизию. Текущие TypeSpec refs имеют обязательные числовые RN; это наблюдение о current state и не ограничивает целевой API.


## 16. Дополнительная проверка ссылок внутри aggregate state

Аудит term variable persistence выявил важное исключение из предположения, что TypeRef всегда сохраняется как полная пара ID + RN:

- В core termvar.VarRec поле TypeRef имеет тип typesem.SemRef и содержит TypeID + TypeRN.
- Persistence DTO internal/adt/termvar/adapter/VarRec содержит только TypeID (JSON key id), ChnlPH и ExpCK.
- Generated converter DataFromVarRec сохраняет TypeID, но не TypeRN; обратный converter восстанавливает TypeID, оставляя TypeRN равным zero-value.
- ExpCK сохраняется отдельно и указывает на конкретное type expression state.

Таким образом, вложенная ссылка переменной на Type сейчас фактически является ссылкой на identity Type без RN, а выбранное expression state фиксируется отдельным ExpCK. Это не просто потеря поля в converter: persistence schema намеренно не содержит TypeRN. Нужно определить, является ли это правильной моделью snapshot semantics (Type identity + pinned expression), либо при переходе на общий Description root переменная должна ссылаться на root RN, а ExpCK останется ссылкой на конкретное expression state.

Следствие: нельзя автоматически объявить все Type/Term references durable revision-bearing references с обязательным RN. Нужно классифицировать каждое поле:
- exact aggregate state reference (ID + RN);
- identity-only association (ID);
- latest selector (ID + zero RN);
- pinned child-state reference (например, ID + ExpCK).

Это особенно важно для ссылочной целостности и concurrency: изменение Type aggregate может увеличить root RN, но не обязательно должно инвалидировать уже сохранённую term variable, если она намеренно закрепляет прежний ExpCK.


## 17. Зафиксированные решения по границе descsem и терминологии

В обсуждении issue #128 подтверждены следующие решения:

1. Использовать термин **QN** (*qualified name*) вместо `QN` во всём ADR и последующей документации.
2. Единый `descsem` aggregate существует **внутри каждого bounded context**. Следовательно, `proc` и `pool` владеют независимыми descriptor identity spaces, QN namespaces, kind и RN. Общая identity space между контекстами не предполагается.
3. Type и Term остаются различными сущностями первого уровня одного `descsem` aggregate в каждом контексте.

Это решение снимает вопрос 15.4(1): root scope — контекстный. Оно не означает общий Go/API repository между `proc` и `pool`; каждый bounded context сохраняет собственные application ports, persistence adapters и транзакционные границы.

## 18. Дополнительный аудит ссылок и TypeExp

### 18.1. TypeExp — адресуемое состояние, а не часть идентичности Type

В обеих реализациях Type root хранит `ExpCK`, а TypeExp state хранится в отдельной context-owned таблице (`proc_type_exps` или `pool_type_exps`). Дерево expression восстанавливается по `exp_ck` через recursive query. В Process Type Modify при изменении expression создаётся новый expression record, обновляется `ExpCK` Type root и увеличивается `TypeRN`; сама expression state адресуется собственным `ExpCK`.

Следовательно, на текущем уровне данных:
- Type identity и текущая версия Type принадлежат Type root;
- `ExpCK` выбирает конкретное expression state;
- expression records не содержат TypeRN и не являются самостоятельными Type roots;
- изменение expression меняет Type root RN и его указатель `ExpCK`, но expression key остаётся отдельным адресом состояния.

Это согласуется с моделью, где TypeExp — kind-specific state внутри Type variant одного Description root, а не самостоятельный descriptor. Однако при переносе в новую схему нужно сохранить возможность закреплять конкретный `ExpCK` независимо от текущего root RN.

### 18.2. Классификация ссылок по проверенным потребителям

| Поле/путь | Core representation | Сохраняемая форма / поведение | Предварительная семантика |
|---|---|---|---|
| Process Type Modify/Retrieve | Type ID + Type RN | Root row хранит ID, RN и ExpCK; Modify сравнивает RN | Exact current-state reference / optimistic concurrency |
| Pool Type Create | Type ID + Type RN | Root row хранит ID, RN и ExpCK | Root identity + current state |
| Pool TermDef Create | Term ID + Term RN | Term root row хранит ID и RN; variable payload отдельно | Root identity + current state |
| Process TermDec Create/Retrieve | Term ID + Term RN | Root row хранит ID и RN; текущая реализация Retrieve проверяет ID без полноценного RN contract | Требует выравнивания с целевым exact/latest правилом |
| TermVar.TypeRef | Type ID + Type RN в core | DTO хранит только TypeID и ExpCK; RN намеренно не сериализуется | Identity association + pinned TypeExp state |
| Pool SpawnSpec2.ProcTermRef | Term ID + Term RN | Входной API принимает пару ID/RN; в проверенном коде это преобразуется в core SpawnSpec2 | Revision-bearing input reference; дальнейшая persistence/use требует отдельной трассировки |
| Pool SpawnSpec2.ProcCompRefs | Computation ID + Computation RN | Входной API принимает пары ID/RN; core передаёт refs в SpawnSpec2 | Revision-bearing input references |

Важно: последняя пара SpawnSpec2 refs проверена на inbound conversion в core, но в этом аудите не установлено, что сами refs сохраняются как durable references. Не следует делать такой вывод без трассировки полного Spawn use case.

### 18.3. Выводы для целевого descsem aggregate

- Description root владеет общей для Type и Term парой identity/RN внутри конкретного контекста.
- Type и Term payload остаются специализированными вариантами; TypeExp остаётся адресуемым child state через `ExpCK`.
- Вложенная связь на Type может быть identity-only и одновременно закреплять конкретный expression state. Поэтому root RN и child-state key решают разные задачи и не должны автоматически сливаться.
- Контекстная независимость означает, что `proc` reference нельзя разрешить через `pool` repository (и наоборот), даже если низкоуровневая пара ID/RN имеет одинаковую форму.
- Входной revision-bearing reference ещё не доказывает durable persistence. Для каждой связи отдельно фиксировать: источник (API/core), сохранение, rehydration и точку проверки RN.

### 18.4. Оставшаяся трассировка перед проектированием DDL

1. Проследить Spawn use case целиком: от `SpawnSpec2` до создания process/computation и выяснить, где проверяются `ProcTermRef` и `ProcCompRefs`, и сохраняются ли они.
2. Завершить инвентаризацию computation references в variables, connections и actions, включая RN check/CAS.
3. Сопоставить все ссылки Type/Term в обоих контекстах с четырьмя классами из раздела 16.
4. На основании этой карты утвердить root-level revision policy и спроектировать новую schema/API migration без обязательства сохранять существующие TypeRef/TermRef wire DTO.

До выполнения этих шагов код и DDL не менять.

## 19. Целевая модель: description aggregate

Этот раздел описывает целевую предметную модель независимо от текущих Go packages и SQL tables. Предыдущие разделы фиксируют current state; при конфликте целевое предложение ниже является проектной моделью. Реализация и DDL не начинаются до согласования открытых решений в §23.

### 19.1. Границы и владение

В каждом bounded context существует собственный descsem aggregate, далее **Description**:

- proc.Description принадлежит только контексту proc;
- pool.Description принадлежит только контексту pool;
- у контекстов независимые пространства identity, QN, kind и revision;
- одинаковые ID в разных контекстах не обозначают один объект;
- proc-ссылка никогда не разрешается через pool repository, и наоборот.

Description — aggregate root и consistency boundary. Он владеет descriptor identity, QN, kind, root revision и специализированным состоянием ровно одной сущности первого уровня: Type либо Term. Computation остаётся отдельным видом агрегата.

Это не общий aggregate между proc и pool и не общий application service. Каждый контекст имеет собственные use cases, ports, транзакции и persistence adapter.

### 19.2. Состав Description

Концептуальная модель:

~~~text
Description
  identity: DescID
  qn: QN
  kind: Type | Term
  entityID: TypeID | TermID
  revision: DescriptionRevision
  lifecycle: active | retired

  state:
    TypeState
      type-specific attributes
      currentTypeExp: ExpCK
    | TermState
      term-specific attributes
      declaration/definition facets
~~~

Инварианты:

1. Description имеет ровно один identity и ровно один текущий kind.
2. Kind выбирается при создании и неизменяем на протяжении жизни Description. Преобразование Type в Term (или наоборот) — не mutation существующего объекта, а создание нового Description с новой identity.
3. Для одного Description существует ровно одно соответствующее specialized state: TypeState для Type либо TermState для Term. Состояние другого kind отсутствует.
4. Type и Term — разные сущности первого уровня и разные типизированные операции/ссылки. Общее владение root не превращает их в один доменный тип и не разрешает взаимозаменять Type и Term.
5. Все изменения Description и его owned state, влияющие на наблюдаемое состояние агрегата, проходят через root и увеличивают одну общую revision.
6. QN уникален в пределах своего bounded context среди активных Description. QN — имя/способ поиска, но не identity. Переименование не меняет identity, однако является mutation Description и увеличивает revision.
7. Retirement должен сохранять невозможность случайно переиспользовать identity. Освобождение QN после retirement — отдельная политика; физическое удаление не является частью базовой модели.

DescID идентифицирует Description aggregate root и не переиспользуется как identity принадлежащей ему доменной сущности. TypeID и TermID — отдельные identity соответствующих сущностей первого уровня. Type и Term сохраняют разные доменные контракты; kind неизменяем и определяет, какая сущность принадлежит Description.

### 19.3. QN — отдельное имя, не identity

QN является уникальным именем Description в пространстве конкретного контекста. Он используется для discovery/resolution и не встраивается в durable identity reference.

- Переименование QN не меняет DescID.
- Сохранённые ссылки продолжают указывать на тот же Description после переименования.
- QN lookup возвращает типизированную ссылку (ID + текущая revision) и kind.
- Уникальность QN обеспечивается атомарно на уровне контекста.
- Разрешение QN и чтение состояния должны быть согласованы в одной транзакционной точке.
- Alias/history для старых QN не вводится автоматически: если понадобится, это отдельная сущность и политика.

### 19.4. Root revision

DescriptionRevision — монотонная версия текущего состояния aggregate root. Она общая для Type и Term, а не отдельная для каждого прежнего Type/Term storage object.

- Создание задаёт начальную ненулевую revision.
- Каждая успешная mutation Description или owned Type/Term state увеличивает revision ровно один раз за транзакцию.
- CAS принимает expected revision и атомарно изменяет root только при совпадении.
- Изменение QN, TypeState, TermState, активного ExpCK или owned child state увеличивает root revision.
- Чтение по identity без проверки revision допускается как отдельный latest-read use case.
- Историческое чтение по старой revision не подразумевается; для него нужна отдельная history/versioned storage.

Revision защищает согласованность Description aggregate, но не заменяет revision других агрегатов. Изменение Computation не увеличивает Description revision и наоборот.

### 19.5. TypeExp как адресуемое дочернее состояние

TypeExp — специализированное дочернее состояние Type, не отдельный Description и не самостоятельная Type identity.

- TypeState содержит ссылку на текущий ExpCK.
- ExpCK идентифицирует конкретный expression state независимо от текущей Description revision.
- Создание нового TypeExp создаёт новый ExpCK; переключение TypeState на него увеличивает Description revision.
- Если другой объект сохраняет пару (DescID, ExpCK), он закрепляет конкретное expression state и не обязан закреплять root revision.
- Старые expression states, на которые существуют durable pins, должны оставаться доступными. Удаление возможно только после гарантированного отсутствия ссылок либо по явной retention policy.
- Формат и схема TypeExp могут различаться между proc и pool; общий Description root не требует унифицировать payload или expression semantics.

Root revision отвечает за актуальность целого Description, ExpCK — за адресацию конкретного дочернего expression state.

### 19.6. Term declaration и Term definition

Term — один вид сущности первого уровня Description. termdec и termdef следует рассматривать как контекстные роли/операции над Term, а не как основания создавать второй независимый descriptor identity space.

Declaration и Definition — разные facets одной Term entity, а не самостоятельные сущности первого уровня и не разные lifecycle-типы. Создание Term создаёт его Declaration; Definition может отсутствовать до отдельной команды создания/заполнения. Обе части принадлежат одному TermState и изменения проходят через Description root revision.

В proc и pool должны действовать одинаковые семантические операции и lifecycle guarantees. Историческое разделение proc/termdec и proc/termdef не задаёт целевую пакетную структуру. Вопрос о названии и границах Go-пакетов (например, term или termdef) отложен и не является решением этого ADR.

### 19.7. Computation aggregate: Execution и Connection

Computation — самостоятельный aggregate root в каждом bounded context, отличный от Description. Внутри Computation находятся две сущности первого уровня: **Execution** и **Connection**.

~~~text
Computation (aggregate root)
  identity: ComputationID
  revision: ComputationRevision

  entities:
    Execution
      execution state
      variables / channel bindings
    Connection
      connection identity
      provider endpoint
      client endpoint
      offset / interaction state
      acts (Connection-owned ordered records)
~~~

Это целевая ownership model; схема не требует, чтобы физическая БД повторяла такое дерево таблицами или FK один-к-одному.

- Execution и Connection — разные entity types первого уровня Computation. Connection не является вложенной частью Execution.
- Execution содержит состояние выполнения Computation и принадлежащие ему переменные/channel bindings. Изменения этого состояния выполняются через Computation root и защищаются его revision.
- Connection имеет собственную identity и связывает Provider и Client Computations. Для одной Connection aggregate owner — Provider Computation; Client Computation является внешним endpoint и не становится вторым владельцем этой Connection.
- Acts — упорядоченные записи взаимодействия, принадлежащие Connection; они не являются самостоятельными first-level entities Computation и не образуют отдельный aggregate root.
- Изменение Connection выполняется через owning Provider Computation и защищается expected Provider Computation revision. Отдельная Connection revision для OCC не вводится.
- Client Computation revision увеличивается только если та же операция меняет её собственный Execution state. Само изменение Connection не увеличивает revision клиентского агрегата.
- Provider/Client endpoint references адресуют Computation aggregates, а не Description.
- CompRef адресует Computation aggregate. Execution и Connection имеют собственные entity identities/keys в пределах определённых для них контрактов; конкретная cardinality и необходимость отдельного ExecutionID должны быть уточнены в §23.
- Изменение нескольких Computations требует явной транзакционной стратегии и набора expected revisions. Нельзя считать, что ссылка на оба endpoint автоматически означает атомарное владение обоими агрегатами.

Таким образом, Computation aggregate имеет одну revision consistency boundary, охватывающую принадлежащие ему Execution и Connection entities. Для Connection, принадлежащей Provider Computation, именно Provider root revision является concurrency token даже когда операцию инициирует Client endpoint.

## 20. Целевые контракты ссылок

Ссылка должна описывать, что адресуется и какую гарантию ожидает операция. Любое поле ID+RN не следует трактовать как один и тот же контракт.

### 20.1. Типизированные identity references

На уровне предметного контракта используются разные семейства:

- proc.TypeRef и proc.CompRef;
- pool.TypeRef и pool.CompRef;
- Term адресуется непосредственно через TermID в соответствующем bounded context.

Контекстная граница задаётся repository/use case, а не отдельным TermRef wrapper. TermID идентифицирует Term entity; DescID идентифицирует владеющий Description root. Эти идентификаторы не взаимозаменяемы. Revision не является частью identity TermID: expected revision передаётся отдельно как предусловие команды, а latest-read выбирается явной семантикой запроса.

TypeRef адресует Type по TypeID и разрешается через соответствующий Type Description. CompRef адресует отдельный Computation aggregate. Старые публичные DTO не являются ограничением этой целевой модели и могут быть заменены вместе с API.

### 20.2. Семантические режимы

| Семантика | Данные | Гарантия |
|---|---|---|
| Identity reference | ID | Указывает на ту же сущность независимо от текущей revision |
| Latest selector | ID + ZeroRevision либо отдельный latest-query | Прочитать актуальное состояние; не является durable snapshot |
| Expected-current reference | ID + nonzero expected revision | Операция разрешена только если текущая revision совпадает |
| Pinned child-state reference | ID + ExpCK | Указывает на конкретное дочернее состояние, не на всю root revision |
| Historical revision reference | ID + revision | Не поддерживается базовой моделью; требует versioned history |

Zero revision допустима только для query/selector semantics. Сохраняемая durable ссылка не должна использовать ZeroRevision как обещание конкретного состояния. Если долговременная связь означает «тот же объект», она хранит identity. Если она означает «именно это состояние», она хранит отдельный pin на immutable/versioned child state или требует versioned aggregate history.

Expected-current reference — не историческая ссылка. Она выражает предусловие CAS/command, а не запрос к старому snapshot.

### 20.3. Правила разрешения ссылок

Typed resolver выполняет проверки:

1. Scope: искать только в repository соответствующего bounded context.
2. Existence: root identity существует и доступна.
3. Kind: TypeRef разрешается только в Type Description; TermID разрешается только в Term Description.
4. Revision policy: latest read игнорирует expected revision только если use case явно выбрал latest semantics; command с expected revision проверяет exact equality.
5. Child pin: ExpCK принадлежит указанному Type Description и существует.
6. Lifecycle: retired/deleted state не разрешается обычными active-use cases.

Нельзя молча трактовать ненулевую неправильную revision как latest, а нулевую — как подтверждение текущего состояния.

### 20.4. Durable associations и проверка revision

Каждое сохраняемое поле ссылки классифицируется по смыслу:

- **Identity association** — сохраняется только ID; изменение revision цели не инвалидирует связь.
- **Pinned child association** — сохраняются ID + child key (например, ExpCK); гарантируется доступность именно child state.
- **Version-checked command input** — ID + expected revision передаются команде, но не обязательно сохраняются.
- **Snapshot association** — ID + immutable snapshot/version identity; требует явной versioned model.
- **Latest lookup** — QN или ID используется для чтения актуального состояния, без обещания неизменности после чтения.

Persisted reference не должна автоматически CAS-проверяться при каждом чтении. Проверка нужна только если доменный контракт связи требует именно эту гарантию. Наличие RN во входном API также не доказывает, что его следует сохранять в durable state.

## 21. Target ports и use-case boundaries

Порты принадлежат core соответствующего bounded context. descsem определяет операции над Description aggregate, но не становится универсальным CRUD repository для всех доменных сущностей.

Концептуальные операции Description:

- CreateType / CreateTerm: создать root, QN, kind и соответствующее specialized state атомарно;
- ResolveTypeByQN / ResolveTermByQN: получить типизированную ссылку и согласованное состояние по QN;
- GetType / GetTerm: latest read по identity;
- ModifyType / ModifyTerm: command с expected revision, одна root mutation и одна новая revision;
- RenameDescription: изменить QN с проверкой уникальности и expected revision;
- RetireDescription: завершить жизненный цикл по явной политике;
- Type-specific operations управляют TypeExp и ExpCK через Type aggregate boundary.

Это логические use cases, а не требование реализовать один интерфейс со всеми методами. Разные операции могут быть разнесены по capability packages. Инварианты: root owns consistency, context owns ports, adapters implement them, transaction boundary covers root plus specialized state.

Term declaration/definition ports остаются раздельными там, где различаются команды или семантика, но используют общий Term identity/root contract. Computation ports остаются отдельными.

## 22. Persistence и транзакционные инварианты целевой модели

Целевая схема не обязана повторять текущую форму таблиц. Она должна обеспечивать:

1. единый Description identity space Type + Term внутри каждого контекста;
2. уникальность QN внутри контекста;
3. однозначное соответствие kind специализированному state;
4. невозможность TypeState для Term root и TermState для Type root;
5. атомарную запись Description root, QN, kind, revision и initial specialized state;
6. атомарное увеличение Description revision вместе с изменением owned Type/Term state;
7. проверку expected revision в том же атомарном write, что и mutation;
8. referential integrity для ID + ExpCK pins и retention expression state;
9. отсутствие FK/lookup через границы proc и pool;
10. независимые identity/revision spaces Computation и Description;
11. Computation root revision защищает изменения принадлежащих ему Execution и Connection;
12. каждая Connection имеет одного aggregate owner (Provider Computation), а Client является внешним endpoint;
13. Act stream принадлежит Connection и упорядочивается в рамках Connection, но не вводит отдельную aggregate revision;
14. переменные/channel bindings принадлежат Execution; их изменение выполняется в consistency boundary Computation;
15. создание/изменение Connection и принадлежащих ей Acts атомарно с CAS owning Provider Computation root;
16. изменение Connection не увеличивает Client Computation revision, если собственный Execution state клиента не меняется.

Физическая реализация может использовать отдельные root/entity tables либо иной storage design, если constraints и транзакции обеспечивают те же инварианты. Наличие единой таблицы само по себе не является целью. Нельзя выбирать DDL до выбора точных lifecycle/cardinality правил и mapping всех durable references.

## 23. Открытые предметные решения перед утверждением модели

1. **QN rename/retirement:** подтверждаем ли mutable QN и уникальность только среди active Description; нужен ли alias/history?
2. **Term facets:** могут ли declaration и definition сосуществовать у одного Term, или это взаимоисключающие lifecycle states? Может ли Term существовать до definition?
3. **Term payload:** какие части данных принадлежат declaration, какие definition, какие общие TermState?
4. **TypeExp retention:** обязаны ли все durable ExpCK pins сохраняться бессрочно, либо допустим управляемый retention?
5. **Description revision:** должны ли все metadata изменения (включая rename/lifecycle) увеличивать ту же revision, или нужна отдельная metadata revision?
6. **Execution identity/cardinality — решено:** одна Computation содержит ровно один Execution; ComputationID является его идентификатором. Отдельный ExecutionID и несколько Execution/retries под одной Computation не вводятся. Повторный запуск трактуется как изменение состояния этого Execution.
7. **Reference durability:** для каждого поля Type/Term/Comp refs в Execution variables, Connection endpoints, Acts и Spawn определить identity association vs pinned state vs command precondition. Полный Spawn use case ещё не подтверждён аудитом: SpawnSpec2 доказывает только входную передачу revision-bearing refs.
8. **Cross-aggregate atomicity:** какие команды могут атомарно менять несколько Computations и каков набор expected revisions/lock ordering?
9. **Deletion:** допускается ли физическое удаление Description/Term/Type/Computation и какие durable references блокируют его?

## 24. План перехода (не часть target model)

План миграции будет уточнён после утверждения §19–23:

1. Зафиксировать предметные решения и полную карту durable references.
2. Спроектировать target logical schema и constraints без наследования текущих имён/таблиц.
3. Спроектировать migration/backfill: объединение Type/Term identity space внутри каждого контекста, перенос QN/kind, root revision, specialized states и pinned ExpCK references; определить обработку конфликтов ID/QN и неконсистентных записей.
4. Спроектировать новый набор core ports/use cases и context-specific adapters, не сохраняя старые Go types ради совместимости.
5. Перестроить inbound/outbound API contracts по целевой модели; совместимость со старым публичным API не является требованием этой задачи.
6. Перенести persistence, добавить invariant tests, CAS/concurrency tests, resolver/kind tests и migration verification tests.
7. Удалить старые adapters/contracts после перехода всех consumers.

До отдельного согласования target model и migration plan Go code и DDL остаются неизменными.


## 25. Уточнение модели Computation относительно ADR-0002

ADR-0002 фиксирует важные правила Connection lifecycle и optimistic concurrency, но использует формулировку «Connection — вложенная сущность агрегата Provider Computation». Для целевой модели issue #128 эта формулировка уточняется:

- Computation остаётся aggregate root.
- Execution и Connection — разные сущности первого уровня внутри Computation aggregate.
- Connection не вложена в Execution; её aggregate owner — Provider Computation.
- Acts принадлежат Connection как её упорядоченное внутреннее состояние.
- Правило ADR-0002 о CAS по Provider Computation revision сохраняется: это root-level concurrency control, а не собственная revision Connection.
- Client endpoint не становится владельцем Connection; client root revision меняется только при изменении собственного Execution state.

Иными словами, «сущность первого уровня Computation» описывает место Connection в доменной модели, а «Provider Computation владеет Connection» определяет конкретный aggregate owner для каждой Connection. Эти утверждения совместимы.

## 26. Дополнительный аудит Computation и Spawn

### 26.1. Фактическое представление Computation

В текущей реализации pool и proc имеют отдельные context-owned Computation repositories и storage:

- *_comp_execs хранит comp_id, comp_rn и liab_mode;
- *_comp_vars и специализированные *_struct_vars / *_linear_vars хранят состояние переменных и channel bindings, адресованное через comp_id и comp_rn;
- *_comp_conns хранит conn_id, Provider и Client Computation refs, offset_nr и act_nr;
- *_conn_acts хранит ordered Act records с conn_id, act_nr, comp_id, comp_rn, chnl_id, kind и expression payload.

Это подтверждает фактическое разделение execution state, Connection и Act stream. Целевая модель §19.7 задаёт доменную ownership: Execution и Connection — сущности первого уровня Computation; Act принадлежит Connection. Текущие названия comp_execs и compexec не следует автоматически принимать за окончательную доменную терминологию.

### 26.2. Connection ownership и OCC

ADR-0002 уже задаёт правило: Connection имеет одного owner — Provider Computation. Изменение Connection защищается CAS по ожидаемой Provider Computation revision; отдельная Connection revision не используется. Client Computation revision увеличивается только при изменении её собственного состояния.

Текущий Pool use case Take собирает все затронутые Computation refs, выполняет CAS по ним до записи Connection/variable state и добавляет Provider ref для Connection mutation. Process path также передаёт Provider ref в ExecMod.CompRefs; его repository обновляет соответствующие proc_comp_execs по ID + expected RN. Это поддерживает target rule о Computation root revision как consistency/concurrency boundary для owned Execution и Connection.

### 26.3. Spawn refs: фактический статус

Текущий inbound adapter pool/compexec/adapter/me_echo.go преобразует PoolTermSpawnSpec:

- proc_dr → Term reference с ID + RN;
- каждый элемент proc_irs → Computation reference с ID + RN;
- результат помещается в termexp.SpawnSpec2{ProcTermRef, ProcCompRefs}.

Однако аудит использования SpawnSpec2 показывает, что в текущем core он только объявлен и создаётся inbound adapter-ом; consumer в pool/compexec не найден. pool/compexec.service.Spawn принимает StepSpec, но текущая реализация создаёт новый Process Computation execution record и не использует переданное PoolExp. Поэтому нельзя утверждать, что ProcTermRef/ProcCompRefs разрешаются, проверяются, сохраняются или rehydrate-ятся как durable references. Это незавершённый use case, а не основание выводить контракт хранения из текущего кода.

Следствия для target design:

1. Spawn должен стать полноценным application use case с явным разрешением Term Description и referenced Computations в Process context.
2. Нужно определить, какие Computation refs являются входными expected-current preconditions, какие — identity associations, и какие должны переживать Spawn как durable associations.
3. Нужно определить, должен ли Spawn создавать новую Execution entity в целевом Computation aggregate, новую Computation identity или оба объекта; текущий код недостаточен для такого решения.
4. Нельзя переносить входные RN в постоянное хранилище только потому, что они присутствуют в inbound DTO.
5. Spawn semantics и reference durability должны быть согласованы с ownership и CAS rules Computation из §19.7 и §20.



## 27. Аудит Execution и durable Computation references

### 27.1. Execution identity и cardinality в текущем коде

В обоих контекстах текущая persistence-модель фактически связывает одну execution-запись с Computation identity:

- `proc/compexec.ExecRec` содержит только `CompRef` и `LiabMode`; отдельного ExecutionID в core record нет.
- `pool/compexec.ExecRec` имеет ту же форму. `Run` создаёт новый `CompRef`, затем атомарно записывает computation execution row и начальные variables.
- Оба Ent adapters создают запись `*_comp_execs` с `comp_id`, `comp_rn`, `liab_mode`; отдельный ключ Execution в рассмотренных core/adapter contracts не обнаружен.
- Variables/channel bindings адресуются через Computation ID и revision, а не через самостоятельный ExecutionID.

Следовательно, текущая реализация моделирует один execution-state record на Computation identity. В целевой модели принимается то же решение: **одна Computation содержит ровно одну сущность Execution, а ComputationID служит её идентификатором**. Отдельный ExecutionID не вводится. Повторный запуск в рамках этой модели изменяет состояние существующего Execution; создание нескольких запусков или истории retries под одной Computation не входит в текущую модель и потребует отдельного пересмотра.

### 27.2. Variables и channel bindings

Общий ADT `compvar` задаёт `StructRec` и `LinearRec` с полями `CompRef`, `ConnID`, `ChnlID`, `ChnlPH`, `ChnlBS`, `ExpCK`. Обе разновидности представляют состояние канала Computation; `ExpCK` указывает на expression state типа канала.

В Pool execution snapshot загружает struct/linear variables по `comp_id`; в Process execution snapshot загружает linear variables для Computation. Process mutation формирует набор затронутых `CompRef` и CAS-обновляет соответствующие `proc_comp_execs` перед записью bindings. Pool mutation аналогично собирает affected Computation refs, CAS-ит root revisions и затем записывает variable state.

Текущая модель тем самым привязывает variable/channel state к Computation execution state и защищает её через Computation revision. Но само наличие `CompRef` в каждой записи не доказывает, что RN целевого Computation должен быть частью долговременной association: в target следует отдельно решить, является ли поле ссылкой на владельца (тогда достаточно owner identity), либо исторической provenance/expected-current token.

### 27.3. Connection endpoints

В Pool `ConnRec` содержит Provider и Client `CompRef` (оба с ID + RN). Ent adapter сохраняет обе пары в `pool_comp_conns`. При чтении Connection snapshot запрос получает актуальный Provider RN через join с `pool_comp_execs`; сохранённый endpoint RN не используется как immutable historical snapshot в этом read path.

В Process Connection adapter сохраняет Connection state и при чтении получает Provider ID + актуальный Provider RN через join с `proc_comp_execs`. В просмотренном Process core contract `ConnMod/ConnSnap` содержит Provider ID + RN; отдельный Client endpoint ref в этих структурах не представлен.

Вывод для target model:

- Provider и Client — endpoint associations на Computation identities.
- Provider — единственный aggregate owner Connection; его expected-current revision используется для CAS изменения Connection.
- Client endpoint не становится co-owner и его RN не должен автоматически превращаться в CAS precondition для каждой Connection mutation.
- Если требуется зафиксировать историческое состояние endpoint на момент создания Connection, это отдельная snapshot/provenance семантика, которую текущий read path не подтверждает.

### 27.4. Acts: participant reference и ordering

Оба контекста хранят Act records с `conn_id`, `act_nr`, `comp_id`, `comp_rn`, `chnl_id`, kind и expression payload. Adapter присваивает Act номер, равный следующей Provider Computation revision (`ProviderCompRN + 1)); чтение Connection упорядочивает записи по `act_nr` и выбирает записи после `offset_nr`.

В core Act record `CompRef` указывает на Computation, чьё действие/продолжение представлено записью. Поэтому `comp_id + comp_rn` в Act сейчас выглядит как participant/provenance stamp, а не как ссылка-владелец Act и не как доказательство, что историческую revision Computation можно восстановить. Текущая модель не предоставляет historical Computation snapshots.

Target contract должен развести:
- Connection ownership and ordering: `ConnectionID + ActOrder` внутри Provider Computation aggregate;
- participant identity: Computation ID, если достаточно знать участника;
- historical provenance: отдельный immutable event/snapshot identifier, если необходимо доказать состояние участника на момент Act;
- CAS: expected Provider Computation revision на команду append/update, не durable RN участника.

Не следует сохранять Act participant RN как будто оно гарантирует возможность исторического чтения, пока не введена versioned Computation history.

### 27.5. Сводная классификация текущих полей и целевых решений

| Поле / место | Что подтверждает current code | Предварительный target meaning |
|---|---|---|
| Variable `CompRef` | Computation, чьё variable/channel state записывается | Owner identity; RN — только command precondition, если явно требуется |
| Variable `ExpCK` | Конкретное TypeExp state канала | Pinned child-state reference к Type Description + ExpCK |
| Connection Provider ref | Owner Connection; Provider RN используется для root CAS | Provider Computation identity + expected-current RN на mutation |
| Connection Client ref | Второй endpoint; не owner Connection | Client Computation identity; RN не CAS-ить автоматически |
| Act `CompRef` | Computation participant, связанный с Act | Participant identity; историческую provenance моделировать отдельно при необходимости |
| Act `CompRN` | Revision stamp, записанный при создании Act | Не считать historical snapshot reference без versioned Computation storage |
| Spawn `ProcTermRef` | Передаётся в не завершённом inbound flow | Term Description identity + отдельно определённая command revision policy |
| Spawn `ProcCompRefs` | Передаются в не завершённом inbound flow | Для каждого определить identity association или expected-current precondition; не сохранять RN автоматически |

Эта таблица разделяет наблюдаемую текущую форму данных и target semantics. Последняя колонка — проектная рекомендация, а не утверждение о существующем поведении.

### 27.6. Зафиксированное решение по Execution identity

Целевая модель принимает одну сущность Execution на Computation:

- ComputationID одновременно адресует Computation aggregate и его единственный Execution entity;
- отдельный ExecutionID не вводится;
- variables/channel bindings адресуются через ComputationID и принадлежат этому Execution;
- повторный запуск в рамках текущей модели изменяет состояние существующего Execution;
- несколько Execution под одной Computation, параллельные запуски и история retries не поддерживаются этим решением и потребуют отдельного изменения модели.

Это решение соответствует наблюдаемой текущей форме core/storage contracts и не требует вводить дополнительную identity без подтверждённой предметной потребности.

Этот аудит завершает первичную карту известных Computation reference sites, но не заменяет полный поиск по всем persistence schemas и inbound/outbound API contracts перед миграцией.

## 28. Аудит Term: declaration и definition

### 28.1. Что фактически представляют текущие данные

В Process core `termdec.DecRec/DecSnap` содержат TermRef и два набора channel variables: `LiabVar` и `AssetVars`. Каждый variable содержит TypeRef и channel placeholder. Это декларационная/интерфейсная часть Term: какие каналы предоставляет или ожидает процесс и как они типизированы.

В TypeSpec отдельно описаны:
- `Proc.DecSpec/DecSnap` и `Pool.DecSpec/DecSnap` для декларационных данных;
- `Proc.Term.ExprSpec` и `Pool.Term.ExprSpec` для выражения определения;
- GET и PATCH `/processes/{id}/definition` и `/pools/{id}/definition` для чтения и аддитивного расширения определения.

README TypeSpec описывает PATCH как атомарное аддитивное расширение: существующие expression parts нельзя удалить или заменить, успешное расширение создаёт новую Term revision. Это контрактная модель TypeSpec; сам по себе он не доказывает полноту текущего runtime implementation.

Pool `termdef.DefRec` содержит те же `LiabVar` и `AssetVars`, что Process `termdec.DecRec`. Его `Create` разрешает Type QNs, создаёт TermRef и сохраняет declaration payload. Отдельного definition payload в этой записи нет.

### 28.2. Incept и незавершённый proc/termdef

Process `termdec.Incept` создаёт новую запись с TermRef и QN, но без LiabVar/AssetVars. Это не переход существующего Term в Definition: последующий `termdec.Create` создаёт новую identity и новую запись, а не модифицирует результат Incept.

Incept используется server-rendered presenter-ом `POST /ssr/decs`. В просмотренном core API отсутствует операция, которая дополняет именно ранее созданную Incept-запись declaration payload-ом. Более того, текущий Process `termdec` adapter сохраняет Term ID/revision и variable payload, но не сохраняет TermQN в `proc_desc_binds`; converter явно игнорирует TermQN. Отдельные SQL запросы adapter-а к `dec_roots` и `desc_sems` не соответствуют таблицам initial migration и не подтверждают работающий root lifecycle.

Process `termdef` не является реализованной второй фазой:
- Create и Retrieve в core возвращают пустые результаты;
- persistence adapter — заглушка;
- пакет не образует работающий definition lifecycle;
- в DI wiring он присутствует, но не подключён как полноценная операция публичного Term API.

Следовательно, current code не доказывает модель lifecycle «Declaration затем Definition». Incept — отдельный незавершённый путь создания пустой записи, а не надёжное основание для целевой state machine.

### 28.3. Целевая рекомендация

Внутри каждого контекста один Term Description root должен владеть разными, семантически типизированными facets:
- **Declaration** — контракт каналов и типов Term (Liab/Asset variables и связанные Type references);
- **Definition** — поведенческое выражение Term и его контекстно-специфичный expression payload.

Declaration и Definition не являются разными видами Description и не требуют разных Term identities. Изменение любого facet, если оно меняет наблюдаемое состояние Term, увеличивает общую Description revision один раз за команду. Формат и операции facets могут различаться между `proc` и `pool`; единый root не требует единого expression payload.

Это рекомендация целевой модели, а не решение о полноте Term state. Открытыми остаются:
1. Может ли Term существовать только с Declaration, только с Definition или с обоими facets; какие из них обязательны при создании?
2. Может ли Declaration быть изменена после создания или она immutable после первоначального определения?
3. Сохраняется ли аддитивная семантика Definition из TypeSpec или целевая модель допускает замену/редактирование выражения?
4. Нужен ли отдельный draft/inception state; если да, это должно быть явное состояние того же Term root с операцией продолжения того же ID, а не создание новой identity через второй Create.

До решения этих вопросов не следует переносить текущий Incept как самостоятельный aggregate type или трактовать proc/termdef как готовую реализацию.

## 29. Аудит QN lifecycle и Description registry

### 29.1. Наблюдаемое хранение и lookup

В initial PostgreSQL migration `proc_desc_binds` и `pool_desc_binds` имеют одинаковую форму:
- `desc_qn ltree UNIQUE`;
- `desc_id varchar`;
- `kind smallint`.

Для QN создаётся GiST index. В этих таблицах нет revision, lifecycle/status, alias/history, FK на Type/Term root или constraint, проверяющего соответствие `kind` фактическому типу root. Уникальность QN действует на всю таблицу; отдельного понятия active/inactive в текущей схеме нет.

Текущий `descsem.Repo` предоставляет только `AddRec`. Adapter выполняет INSERT; операций lookup, rename, retire, delete или истории имён в этом port нет. Type repositories выполняют QN lookup через join с context-specific `desc_binds`. Pool Term repository также разрешает Term по QN через join. Process Term path не имеет согласованной рабочей записи в `proc_desc_binds`: Incept и Create сохраняют Term record, но не добавляют descriptor binding.

В TypeSpec нет операций rename/retire Description. В просмотренных core, adapters, DI wiring и initial migration не обнаружен завершённый QN rename/retirement lifecycle. Это означает отсутствие текущего контракта, а не запрет на такие операции в target model.

### 29.2. Следствия для target model

QN — context-local qualified name и lookup key, но не identity. В target Description root владеет текущим QN, а Type/Term kind остаётся неизменяемым. Resolver должен разрешать QN внутри соответствующего context и возвращать typed Type/Term reference; изменение QN не должно менять DescID.

Рекомендуемая политика:
- QN mutable через явную RenameDescription command;
- rename выполняется атомарно с проверкой уникальности в пределах bounded context и увеличивает общую Description revision;
- uniqueness относится к активным Description roots;
- alias/history не вводить в базовую модель без подтверждённого use case; если позже потребуется стабильное разрешение старых имён, это отдельная сущность/политика, а не неявная часть identity;
- retirement должен освобождать QN только при явно определённой политике ссылочной целостности и повторного использования имени. До решения deletion/retention нельзя автоматически приравнивать retirement к физическому удалению binding.

Это проектная рекомендация. Пользовательское решение по rename, alias/history, retirement и повторному использованию QN ещё не зафиксировано. Целевая схема должна обеспечивать уникальность и атомарность на уровне одного Description aggregate, но конкретная стратегия хранения retired names зависит от согласованной политики.

## 30. Уточнение identity Description и унификация Term lifecycle между контекстами

Это уточнение конкретизирует и заменяет те формулировки §19–21 и §28.3, в которых Description identity и Type/Term identity могли быть прочитаны как один и тот же идентификатор.

### 30.1. Две identity на разных уровнях

В целевой модели необходимо различать:
- **DescID** — identity Description aggregate root в конкретном bounded context;
- **TermID** — identity принадлежащей ему доменной сущности Term.

DescID и TermID — разные идентификаторы с разными назначениями. Description root владеет Term entity, но не переиспользует его ID. Descriptor registry связывает QN и kind с DescID; Term resolver связывает TermID с соответствующим Description root и проверяет, что его kind — Term.

Term адресуется непосредственно по TermID. Отдельный TermRef wrapper не вводится: он не добавляет необходимой семантики поверх TermID. Revision также не является частью TermID и передаётся отдельно там, где команда требует optimistic concurrency. При разрешении TermID система однозначно находит владеющий DescID в пределах текущего context, затем применяет правила kind и lifecycle. Это отображение обеспечивается инвариантами и атомарностью persistence, а не совпадением значений ID.

Аналогичное разделение identity следует применять к Type: TypeID идентифицирует Type entity, DescID — владеющий Description root. Сохранение отдельного TypeRef пока остаётся частью типизированного Type-контракта. Для Computation отдельной Description identity нет: Computation остаётся самостоятельным aggregate root с ComputationID.

DescID не обязан присутствовать в каждой внешней ссылке, если resolver способен надёжно разрешить TypeID/TermID в root. Если конкретному внутреннему use case нужен root identity, он использует DescID непосредственно, не подменяя им identity доменной сущности.

### 30.2. Единый Term lifecycle в proc и pool

Целевое поведение Term должно быть одинаковым в proc и pool. Текущая асимметрия (proc/termdec с Incept/Create и pool/termdef с Create) — недостаток/историческое состояние реализации, а не целевое различие предметной модели.

Оба контекста должны предоставлять один и тот же набор семантических операций и одинаковые гарантии:
- создание Term и его Declaration;
- получение Term и его текущего состояния;
- создание Definition, если её ещё нет;
- расширение/изменение Definition согласно общей согласованной политике;
- одинаковая optimistic concurrency policy через revision содержащего Description aggregate;
- одинаковые правила QN resolution, typed references, lifecycle и ошибок.

Контекстно-специфичными могут оставаться expression payload, структуры channel variables, вычислительная семантика и persistence adapters. Это не должно приводить к различному смыслу операций или к разным lifecycle state machines.

Базовая целевая композиция Term Description:
- один Description aggregate root с DescID, QN, Term kind и общей revision;
- одна принадлежащая ему Term entity с отдельным TermID;
- обязательная Declaration facet, создаваемая вместе с Term;
- необязательная Definition facet, создаваемая отдельной командой;
- Declaration и Definition являются частями одного TermState, а не отдельными агрегатами.

### 30.3. Открытые решения lifecycle

Зафиксированы следующие lifecycle decisions:
1. Term создаётся сразу с Declaration; создание пустого Term без Declaration не входит в базовую модель.
2. Definition необязательна при создании и может быть добавлена отдельной командой.
3. Declaration и Definition принадлежат одной Term entity и одному Description aggregate; изменения увеличивают общую Description revision.
4. Семантика изменения Definition (только additive или также replace/edit) и mutability Declaration остаются открытыми; выбранные правила должны быть одинаковыми в proc и pool.
5. Отдельный draft/inception state не вводится в базовую модель. Если появится подтверждённая потребность, он должен быть состоянием той же пары DescID/TermID, а не новой identity.

Incept не переносится как самостоятельная фаза lifecycle. Целевые операции и гарантии должны быть симметричны в обоих контекстах. Пакетная структура реализации остаётся за рамками этого решения.

## 31. Целевые решения по Type expression и QN

### 31.1. Type identity и expression state

Type является доменной сущностью внутри Description aggregate, а не отдельным aggregate root. Его identity отделена от identity владельца:
- `DescID` идентифицирует Description aggregate root;
- `TypeID` идентифицирует принадлежащую ему Type entity;
- `QN` является человеко-понятным внешним идентификатором Description, используемым в API, конфигурациях и поиске.

Type сохраняет ссылку на текущее expression state через `ExpCK`. Изменение Type создаёт новое expression state и переключает текущий `ExpCK`; конкурентные изменения Description защищаются общей Description revision. Отдельная Type revision в целевой модели не вводится.

Полная история версий Type не является частью базовой модели. Старые expression states сохраняются только пока на них существуют действующие ссылки или иные явно определённые зависимости. Удаление/сборка таких состояний должно учитывать referential integrity; конкретный механизм retention относится к persistence design.

### 31.2. QN как внешний идентификатор

QN — человеко-понятный внешний идентификатор Description, а не внутренняя identity. `DescID` и принадлежащий Description `TypeID` или `TermID` остаются стабильными при переименовании QN.

QN изменяем: переименование не создаёт новый Description и не меняет identity его сущностей. В базовую модель не входят aliases и история прежних QN; их можно добавить при появлении конкретного требования к разрешению старых имён или обратной совместимости конфигураций.

Уникальность QN и атомарность переименования должны обеспечиваться в пределах соответствующего bounded context. Конкретные правила retirement и повторного использования QN остаются отдельным решением и не выводятся автоматически из возможности rename.

### 31.3. Граница архитектурного решения

Проверка соответствия запрошенного вида сущности фактическому kind при разрешении `TypeID`/`TermID`, а также детали ошибок resolver-а являются реализационными контрактами. Они не вводят дополнительные identity или aggregate boundaries и не рассматриваются здесь как отдельное архитектурное решение.

## 32. Description revision и typed entity references

Целевая модель использует одну revision на Description aggregate:
- `DescID` — identity Description root;
- `DescRN` — revision изменяемого состояния Description;
- `TypeID` и `TermID` — самостоятельные identity принадлежащих ему доменных сущностей;
- `ExpCK` — identity конкретного Type expression state.

Отдельные `TypeRN` и `TermRN` в целевой модели не вводятся. `TypeRef` в форме пары `TypeID + TypeRN` не сохраняется: Type адресуется непосредственно через `TypeID`. Аналогично Term адресуется непосредственно через `TermID`, без `TermRef` wrapper.

Команды, которым необходима optimistic concurrency, используют ожидаемый `DescRN` владельца Description. При необходимости явного указания aggregate identity команда передаёт `DescID` вместе с expected `DescRN`; эти значения описывают aggregate boundary и concurrency precondition, а не новую entity reference.

Таким образом, entity identity не включает revision. Любое изменение Type или Term, меняющее наблюдаемое состояние Description, увеличивает общий `DescRN` один раз в рамках соответствующей атомарной команды.

