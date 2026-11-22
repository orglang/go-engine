# ADR-0004: Модель корней агрегатов и semantic references

- **Статус:** Proposed (инвентаризация; решения не приняты)
- **Дата:** 2026-11-22
- **Связано:** go-workspace#128, go-engine#188, go-workspace#122

## Цель

Зафиксировать фактическую модель агрегатов в `proc` и `pool`, отделить её от ADT и подготовить обоснованные решения о reference-контрактах и persistence roots. Этот документ пока не разрешает изменения публичных контрактов или схемы БД.

## 1. Оси модели

- **Контекст исполнения:** `proc` или `pool`. Контексты имеют отдельные application/core пакеты, Ent clients и физические таблицы.
- **Вид агрегата:** type, term, computation. Вид определяет бизнес-идентичность и специализированные инварианты.
- **ADT:** набор семантических value types, primitive identity/revision types и semantic ports. ADT не является ещё одним контекстом исполнения и не должен владеть конкретными persistence adapters.

## 2. Предварительная матрица фактически реализованных агрегатов

| Контекст | Type | Term | Computation |
|---|---|---|---|
| `pool` | `pool/typedef` → `pool_type_defs` | `pool/termdef` → `pool_term_decs` | `pool/compexec` → `pool_comp_execs` |
| `proc` | `proc/typedef` → `proc_type_defs` | `proc/termdec` → `proc_term_decs` | `proc/compexec` → `proc_comp_execs` |

`proc/termdef` присутствует в исходниках, но текущий core содержит незавершённые заглушки; в initial SQL migration и Ent schema соответствующего process term-definition root нет. До отдельного подтверждения это не следует считать валидным самостоятельным агрегатом.

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

- `typesem.Repo` предоставляет `TouchRef` и `GetRefsByQNs`; type reference участвует в lookup по QName и optimistic revision flow в Process Type.
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
- Descriptor binding tables (`proc_desc_binds`, `pool_desc_binds`) связывают QName с ID и kind; необходимо отдельно решить, являются ли они частью aggregate consistency boundary или отдельным индексом/registry.

В initial migration нет отдельных root tables поверх type/term/computation tables: эти таблицы сами содержат ID и revision. Это предварительно поддерживает модель «specialized root table is the aggregate root», но не доказывает, что все lifecycle и referential-integrity требования удовлетворены.

## 6. Известные несоответствия, требующие проверки

- `typesem/core.go` содержит `typeKind` и `PoolKind`/`ProcKind`, которые не отражены в текущем reference contract; проверить все usages перед удалением.
- `proc/termdef` выглядит как незавершённый/legacy пакет и использует `typesem.SemRef`, тогда как фактический process term aggregate находится в `proc/termdec` и использует `termsem.SemRef`.
- `proc/termdec/adapter/ds_ent_pgx.go` содержит запросы к `dec_roots` и `desc_sems`; эти таблицы не объявлены в единственной обнаруженной migration `db/postgres/owner/000001_initial.up.sql`. Установить, являются ли запросы устаревшими, схема неполной или объекты управляются вне этого migration tree. Не менять DDL до выяснения.
- В `proc/typedef/adapter/ds_pgx.go` обнаружены SQL-фрагменты со старыми/нестыкующимися именами колонок; проверить, используются ли эти методы текущим DI wiring и актуальны ли запросы.

## 7. Решения, ожидающие завершения аудита

1. Терминология: оставить `SemRef`, перейти к `TypeRef`/`TermRef`/`CompRef` или использовать `AggRef` с типизированными обёртками.
2. Состав reference: ID и revision; роль QName (вероятно, не часть идентичности); инварианты revision для каждого вида.
3. Root model: оставить существующие специализированные root tables или ввести отдельную root table только при доказанной необходимости для целостности/lifecycle.
4. Таблицы aggregate state: для каждой зависимой таблицы подтвердить прямую ссылку на root Ref или обосновать исключение.
5. `typesem` и `termsem`: объединять только если общими являются не только поля, но и бизнес-инварианты/операции.
6. ADT ports: определить, какие semantic ports действительно принадлежат ADT, и не протаскивать Ent/persistence DTO в core.

## 8. Ограничения реализации

- Не менять SDK.
- Не менять SQL migrations без отдельного обоснования и совместимого migration plan.
- Не создавать все декартовы комбинации автоматически.
- Не вводить универсальный CRUD или универсальный reference только ради уменьшения числа типов.
- Не переносить Ent API или persistence DTO за границы adapters.
- Сохранять внешние wire/API contracts до отдельного решения о совместимости.


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

1. Сохранить три семейства семантических ссылок по виду агрегата: Type, Term и Computation. На уровне Go рассмотреть переименование `SemRef` в `TypeRef`/`TermRef`/`CompRef` как отдельный compatibility-aware шаг, а не смешивать его с изменением поведения.
2. Не включать `Proc`/`Pool` в тип ссылки. Контекст остаётся частью владельца, репозитория, application service и persistence boundary.
3. Не объединять Proc/Pool roots, TypeExp state, Connection lifecycle или use cases только из-за общей формы ссылки.
4. Не вводить универсальный `AggRef`: публичная модель уже имеет три конкретных типа, и они точнее выражают допустимые связи.
5. Перед изменением Go contracts проверить все conversion sites, сериализацию, DI wiring и потребителей semantic ports; отдельно определить, должна ли revision быть частью каждой ссылки в каждом use case.


## 13. Отдельная гипотеза: схлопнуть TypeRef и TermRef

Эта гипотеза рассматривается отдельно от вопроса об объединении aggregate roots. У `typesem.SemRef` и `termsem.SemRef` одинаковая физическая форма: пара `identity.ADT` + `seqnum.ADT`, одинаковый конструктор `New()` и одинаковая базовая трактовка revision. Это аргумент в пользу общего value representation.

При этом обнаружены реальные семантические различия:

- Type references разрешаются по QName через `typesem.Repo.GetRefsByQNs`; этот port также предоставляет `TouchRef`. `proc/termdec` использует его, чтобы разрешить типы liability/asset variables.
- Term references сейчас не имеют аналогичного semantic repository port. Term QName разрешается контекстными `termdef`/`termdec` repositories и descriptor bindings.
- `termvar.VarRec.TypeRef` всегда ссылается на Type, тогда как `pool/termexp.SpawnSpec2.ProcTermRef` ссылается на Process Term. В этих местах вид ссылки — часть предметного смысла поля.
- Descriptor binding (`descsem`) различает `TypeKind` и `TermKind`, то есть реестр QName уже хранит явный discriminator вида сущности.
- Публичный TypeSpec API в настоящее время экспортирует `Any.TypeRef` и `Any.TermRef` отдельно. Внутреннее объединение не должно автоматически менять wire contract.

### Варианты модели

1. **Общий внутренний value type, разные именованные поля/контракты.** Ввести общий `adt/semref.Ref` (ID + revision), а Type/Term boundary packages использовать для семантических операций и явных конвертаций. Это устраняет структурное дублирование, но требует не допускать неявной подмены Type ссылкой на Term.
2. **Один полностью общий `SemRef` без discriminator.** Самый компактный вариант, но вызывающий код должен знать вид root из контекста. Он слабее защищает от ошибочного lookup в неправильном repository и не подходит для автономной передачи/хранения ссылки без контекста.
3. **Общий `SemRef` с kind discriminator.** Ссылка становится самодостаточной для маршрутизации, но kind дублирует часть информации, уже заданной типизированным полем/контрактом и descriptor registry. Также нужно определить, входит ли kind в идентичность или только в маршрутизацию.
4. **Оставить TypeRef и TermRef независимыми.** Сохраняет максимальную статическую типобезопасность, но дублирует одинаковую пару ID/revision и конвертеры.

### Предварительная рекомендация

Исследовать вариант 1 как основной: общий низкоуровневый value representation для Type и Term, но сохранять семантические границы в названиях полей, repository ports и use cases. Не объединять `typesem.Repo` и операции поиска по QName с term repositories. Не менять TypeSpec `Any.TypeRef`/`Any.TermRef` в рамках этой задачи, пока не доказана необходимость wire-level изменения.

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

1. Определить границу и состав общего \`descsem\` aggregate: Type, Term, descriptor/QName bindings, их revision и lifecycle.
2. Определить, является ли один \`descsem\` root общим для всего \`proc\`/\`pool\` контекста или каждый контекст имеет собственный root и собственную identity space.
3. Сопоставить ссылки и foreign keys во всех Type/Term consumers с новой aggregate boundary; установить, какие ссылки должны нести RN, а какие могут запрашивать latest.
4. Определить durable reference contract отдельно от read selector contract и выбрать Go representation для optional RN.
5. Включить будущий \`proc/termdef\` в целевую матрицу, не выдавая текущие заглушки за реализованный persistence root.


## 15. Аудит существующего descsem (фактическая реализация)

Проверка исходников показывает, что internal/adt/descsem уже существует. Однако текущий descsem пока **не является полноценным aggregate root для Type и Term**. Это минимальный descriptor/QName binding registry:

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

### 15.2. QName resolution сейчас распределён по потребителям

- Process Type lookup по QName идёт через proc_desc_binds с последующим чтением proc_type_defs.
- Pool Type lookup по QName идёт через pool_desc_binds с последующим чтением pool_type_defs.
- Pool Term lookup по QName идёт через pool_desc_binds с последующим чтением pool_term_decs.
- Process Term lookup в текущем proc/termdec adapter опирается на устаревшие SQL fragments и не является надёжным образцом целевого контракта.
- typesem.Repo.GetRefsByQNs — не общий descriptor resolver: он жёстко соединяет proc_desc_binds с proc_type_defs и возвращает только Type refs.

Иными словами, descsem пока пишет QName/kind binding, но чтение и проверка kind/сущности реализованы в контекстных специализированных repositories. Это создаёт дублирование lookup logic и не обеспечивает единую ссылочную целостность.

### 15.3. Архитектурная гипотеза для целевой модели

Исследовать эволюцию descsem в общий aggregate root внутри каждого bounded context (proc и pool), который владеет descriptor identity, QName, kind и RN. Type и Term в таком случае остаются разными сущностями/вариантами первого уровня одного агрегата; специализированные данные Type (включая TypeExp) и Term хранятся как kind-specific state, принадлежащий этому root.

Это пока гипотеза, а не принятое решение. Необходимо сравнить как минимум две формы:

1. **Единый descriptor root + специализированное состояние.** Один descsem root на descriptor ID, содержащий QName, kind и RN; type/term-specific tables являются state/parts этого root и используют общий ID/RN.
2. **Общий descriptor registry над самостоятельными Type/Term roots.** descsem централизует QName и kind, но Type/Term сохраняют собственные root identities/revisions и consistency boundaries.

Первая форма соответствует сильному прочтению идеи «две сущности первого уровня одного агрегата» и может устранить дублирование identity/revision semantics. Вторая — более эволюционный вариант, но сама по себе не объединяет aggregate boundary и не гарантирует атомарность cross-kind invariants.

### 15.4. Существенные вопросы до выбора схемы

1. Должен ли descsem root существовать отдельно в каждом контексте, как сейчас (proc_desc_binds и pool_desc_binds), или одна descriptor identity должна быть общей между proc и pool? Текущая схема подтверждает только контекстное владение.
2. Является ли QName частью identity root или изменяемым unique name/binding, который можно переименовать независимо от ID?
3. Должен ли RN жить в descriptor root и версионировать изменения Type/Term state совместно?
4. Как моделировать TypeExp: часть Type variant state под тем же root RN или отдельный immutable/versioned child resource?
5. Как обеспечить, чтобы kind соответствовал ровно одному специализированному состоянию и чтобы QName разрешался в согласованный ID/RN?
6. Должны ли Process Term Declaration (termdec) и Process Term Definition (termdef) быть двумя use cases одной Term-сущности/варианта, либо это разные роли/сущности внутри descriptor aggregate?
7. Как мигрировать текущие таблицы и wire refs, не ломая публичные Any.TypeRef и Any.TermRef?

### 15.5. Наблюдение по zero RN

В internal/adt/seqnum уже определён seqnum.Zero, а seqnum.ADT является числовым value type. Это позволяет выразить selector «latest» без pointer/nil: zero RN означает resolve current; ненулевой RN означает проверку соответствия текущей ревизии. Это должно оставаться семантикой query/use case, а не означать, что durable reference может содержать нулевую ревизию. Текущие публичные TypeSpec refs имеют обязательные числовые RN, поэтому wire contract не следует менять автоматически.


## 16. Дополнительная проверка ссылок внутри aggregate state

Аудит term variable persistence выявил важное исключение из предположения, что TypeRef всегда сохраняется как полная пара ID + RN:

- В core termvar.VarRec поле TypeRef имеет тип typesem.SemRef и содержит TypeID + TypeRN.
- Persistence DTO internal/adt/termvar/adapter/VarRec содержит только TypeID (JSON key id), ChnlPH и ExpCK.
- Generated converter DataFromVarRec сохраняет TypeID, но не TypeRN; обратный converter восстанавливает TypeID, оставляя TypeRN равным zero-value.
- ExpCK сохраняется отдельно и указывает на конкретное type expression state.

Таким образом, вложенная ссылка переменной на Type сейчас фактически является ссылкой на identity Type без RN, а выбранное expression state фиксируется отдельным ExpCK. Это не просто потеря поля в converter: persistence schema намеренно не содержит TypeRN. Нужно определить, является ли это правильной моделью snapshot semantics (Type identity + pinned expression), либо при переходе на общий descsem root переменная должна ссылаться на root RN, а ExpCK останется ссылкой на конкретное expression state.

Следствие: нельзя автоматически объявить все Type/Term references durable revision-bearing references с обязательным RN. Нужно классифицировать каждое поле:
- exact aggregate state reference (ID + RN);
- identity-only association (ID);
- latest selector (ID + zero RN);
- pinned child-state reference (например, ID + ExpCK).

Это особенно важно для ссылочной целостности и concurrency: изменение Type aggregate может увеличить root RN, но не обязательно должно инвалидировать уже сохранённую term variable, если она намеренно закрепляет прежний ExpCK.


## 17. Зафиксированные решения по границе descsem и терминологии

В обсуждении issue #128 подтверждены следующие решения:

1. Использовать термин **QN** (*qualified name*) вместо `QName` во всём ADR и последующей документации.
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

Это согласуется с моделью, где TypeExp — kind-specific state внутри Type variant одного descsem root, а не самостоятельный descriptor. Однако при переносе в новую схему нужно сохранить возможность закреплять конкретный `ExpCK` независимо от текущего root RN.

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

- Descriptor root должен владеть общей для Type и Term парой identity/RN внутри конкретного контекста.
- Type и Term payload остаются специализированными вариантами; TypeExp остаётся адресуемым child state через `ExpCK`.
- Вложенная связь на Type может быть identity-only и одновременно закреплять конкретный expression state. Поэтому root RN и child-state key решают разные задачи и не должны автоматически сливаться.
- Контекстная независимость означает, что `proc` reference нельзя разрешить через `pool` repository (и наоборот), даже если низкоуровневая пара ID/RN имеет одинаковую форму.
- Входной revision-bearing reference ещё не доказывает durable persistence. Для каждой связи отдельно фиксировать: источник (API/core), сохранение, rehydration и точку проверки RN.

### 18.4. Оставшаяся трассировка перед проектированием DDL

1. Проследить Spawn use case целиком: от `SpawnSpec2` до создания process/computation и выяснить, где проверяются `ProcTermRef` и `ProcCompRefs`, и сохраняются ли они.
2. Завершить инвентаризацию computation references в variables, connections и actions, включая RN check/CAS.
3. Сопоставить все ссылки Type/Term в обоих контекстах с четырьмя классами из раздела 16.
4. На основании этой карты определить root-level revision policy, затем предложить schema migration и compatibility strategy для существующих TypeRef/TermRef wire DTO.

До выполнения этих шагов код и DDL не менять.
