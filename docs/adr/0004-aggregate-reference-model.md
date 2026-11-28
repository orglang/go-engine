# ADR-0004: Модель Program и Computation aggregates

- **Статус:** Accepted
- **Дата:** 2026-11-22
- **Связано:** go-workspace#128, go-engine#188, go-workspace#122

## Контекст и область решения

Этот ADR фиксирует целевую доменную модель `proc` и `pool`: границы агрегатов, identity, revision, QN, отношения Type/Term и Computation. Текущая реализация и существующая схема БД не ограничивают целевую модель.

`proc` и `pool` — разные bounded contexts. У каждого собственные агрегаты, identity spaces, application services, порты и persistence adapters. Между контекстами нет общих агрегатов и cross-context references.

Изменения SDK не входят в scope. Конкретная схема хранения, миграции и persistence-механизмы этим ADR не определяются.

## 1. Program aggregate

В каждом bounded context существует Program aggregate root. Он объединяет:

- ровно один Provider Type;
- ноль или более Terms, реализующих этот Provider Type.

Type и Term — разные доменные сущности первого уровня со своими identity и контрактами. Они принадлежат одному Program aggregate, но не становятся взаимозаменяемыми.

Допустим Program только с Type, без Terms. Program без Type и Program только с Terms недопустимы. Term всегда реализует Provider Type своего Program.

Program не имеет собственного QN и не имеет доменного `kind`. Descriptor registry не определяет отдельный вид Program.

### Identity

- `ProgID` — identity Program aggregate root.
- `TypeID` — identity Type.
- `TermID` — identity Term.

Identity Program и принадлежащих ему сущностей различны. Type и Term сохраняют свои собственные IDs даже при общем Program root.

### Revision

`ProgRN` — единственная revision наблюдаемого состояния Program aggregate и общая optimistic concurrency boundary для его мутаций.

Одна атомарная команда, меняющая наблюдаемое состояние Program, увеличивает `ProgRN` ровно один раз. Это правило распространяется на:

- изменение Type expression и переключение текущего `ExpCK`;
- изменение Declaration или Definition Term;
- добавление принадлежащей Program сущности;
- переименование Type или Term;
- Retire Program.

Несколько внутренних записей одной команды не означают несколько increments. Технические изменения без наблюдаемого доменного изменения сами по себе не требуют increment.

Отдельные `TypeRN` и `TermRN` в целевой модели не вводятся. Команды с optimistic concurrency используют expected `ProgRN`.

## 2. QN и единое пространство имён

QN принадлежит непосредственно Type или Term. Program собственного QN не имеет.

Внутри каждого bounded context существует одно общее пространство имён для QN всех Type и Term во всех Program. Поэтому:

- QN любого Type должен быть уникален среди QN всех Type и Terms этого контекста;
- QN любого Term должен быть уникален среди QN всех Terms и Types этого контекста;
- одинаковый QN у Type и Term также запрещён;
- `proc` и `pool` имеют независимые пространства имён.

QN — изменяемое человеко-понятное имя для поиска и внешнего обращения, а не identity. Type QN разрешается в `TypeID`, Term QN — в `TermID`. Типизированные resolver contracts сохраняются; единое пространство имён не требует общего неразмеченного reference type.

Переименование выполняется над Type или Term. Оно атомарно проверяет уникальность QN во всём bounded context, не меняет `TypeID`, `TermID` или `ProgID` и увеличивает `ProgRN` владельца ровно один раз.

Aliases и история прежних QN не входят в базовую модель.

## 3. Type expression и отношения Provider/Client

### Type expression

Type хранит ссылку на текущее expression state через `ExpCK`. Изменение Type создаёт новое expression state и переключает текущий `ExpCK`.

Полная история версий Type не входит в базовую модель. Старые expression states сохраняются, пока на них существуют действующие ссылки или иные зависимости. Удаление таких состояний должно соблюдать referential integrity; конкретные retention и garbage-collection правила относятся к последующему lifecycle/persistence design.

### Provider relationship

Provider — роль Type, контракт которого реализует Term.

- Каждый Term имеет ровно один Provider Type.
- Provider Type и реализующие его Terms принадлежат одному Program.
- Один Provider Type может иметь несколько Terms-реализаций.
- Program может быть создан с Provider Type до появления первого Term.
- Изменение Provider Type допустимо, если все реализующие его Terms корректны в итоговом состоянии.
- Одна атомарная команда может изменить Provider Type и любое количество его Terms. Инварианты проверяются по итоговой композиции, а не по промежуточным состояниям.
- Backward compatibility с предыдущим состоянием Provider Type не требуется.

### Client relationship

Client — роль Type, контракт которого Term использует, но не реализует.

- Term может иметь ноль или более Client references.
- Client Type может находиться в другом Program того же bounded context.
- Cross-bounded-context references запрещены.
- Client reference фиксирует конкретное состояние Type, а не `ProgRN` всего Program.
- Последующие изменения Type не перенаправляют существующий Client reference автоматически.
- Один Type может использоваться несколькими Terms как Client.
- Один Term может использовать собственный Provider Type также и в роли Client; Provider и Client остаются разными отношениями.
- Client reference не передаёт владение Type и не объединяет Program aggregates.

Идентификатор Type state, на который ссылается Client, пока не выбран. `ExpCK` — кандидат, но не подтверждено, что он идентифицирует весь Type contract.

Открытыми остаются: возможность ссылаться на non-current Type state и правила его retention; явное обновление Client reference на новое состояние; возможность нескольких Client uses одного Type с разными pinned states.

## 4. Program lifecycle

Retire применяется ко всему Program aggregate:

- Program становится недоступен для нового использования.
- `ProgID`, `TypeID`, `TermID`, данные и QN принадлежащих Type/Terms сохраняются.
- Retire является наблюдаемой мутацией и увеличивает `ProgRN` один раз.
- QN остаются зарезервированными и не могут быть использованы повторно.

Физический Delete не является частью обычного lifecycle. Он допустим только при отсутствии зависимостей, нарушающих referential integrity. После физического удаления Program его QN могут быть использованы повторно.

## 5. Computation aggregate

Computation — отдельный aggregate root в каждом bounded context. Он не входит в Program.

### Identity и Execution

- `CompID` идентифицирует Computation aggregate и его единственный Execution.
- В одном Computation существует ровно один Execution; отдельный `ExecutionID` не вводится.
- `CompRN` — revision Computation aggregate.
- Variables и channel bindings принадлежат Execution и адресуются через `CompID`.
- Повторный запуск изменяет состояние того же Execution. Несколько Execution, retry history и run history базовой моделью не поддерживаются.

### Connection

Connection — сущность первого уровня внутри Computation aggregate, а не отдельный aggregate root и не часть Execution.

- Connection имеет собственный ID.
- Connection связывает Provider Computation и Client Computation.
- Connection принадлежит Provider Computation; Client endpoint является внешней ссылкой, а не совладельцем Connection.
- Ordered Acts принадлежат Connection и являются его записями, а не отдельными сущностями первого уровня.
- Мутации Connection выполняются через Provider Computation и используют expected revision Provider.
- У Connection нет независимой optimistic concurrency revision.
- Изменение только Connection не увеличивает revision Client Computation. Revision Client увеличивается, только если меняется его Execution.
- Изменения, затрагивающие несколько Computations, требуют явно заданной транзакционной стратегии и expected revisions соответствующих агрегатов. Наличие endpoint reference само по себе не означает совместного владения.

## 6. Архитектурные границы

- Core каждого bounded context владеет своими портами; adapters реализуют их.
- Зависимости направлены внутрь. Domain/application code не зависит от драйверов и конкретных adapters.
- Program operations выполняются в границе Program aggregate; Computation operations — в границе Computation aggregate.
- Type и Term могут иметь раздельные use cases и capability packages, сохраняя общий Program consistency boundary.
- SDK и публичные wire contracts не являются ограничением для целевой внутренней модели.
- Persistence schema и migration strategy проектируются отдельно и должны обеспечивать описанные identity, uniqueness, atomicity и referential-integrity guarantees.

## 7. Отложенные вопросы

Этот ADR намеренно не определяет:

- конкретную физическую схему хранения Program, Type, Term, QN registry и Computation;
- миграцию существующих descriptor bindings к QN → TypeID/TermID;
- механизм обеспечения общей QN uniqueness constraint;
- точный идентификатор и retention policy pinned Type states для Client references;
- правила чтения retired Program, возможного восстановления и физического удаления;
- транзакционную реализацию операций, изменяющих несколько Computation aggregates.

Эти вопросы не меняют зафиксированные выше доменные границы и инварианты.
