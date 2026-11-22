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
