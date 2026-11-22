# OrgLang Engine

Рантайм исполнения OrgLang на Go.

## Структура

- `cmd/engine/` — точка входа и задачи сборки/запуска исполняемого приложения.
- `internal/app/` — прикладные сценарии и входящие адаптеры (HTTP и другие интерфейсы).
- `internal/adt/`, `internal/pool/`, `internal/proc/` — внутренние доменные возможности движка, сгруппированные по предметным областям.
- `internal/lib/` — общие технические контракты и инфраструктурные компоненты, используемые ядром.
- `db/` — инфраструктура БД: миграции, Compose и задачи локального обслуживания.
- `internal/proto/` — внутренние протоколы и версии форматов/исполнения.
- `test/` — сквозные проверки, архитектурные проверки и тестовые сценарии.

`internal/` содержит реализацию движка, не предназначенную для импорта внешними модулями. `pkg/` предназначен для публичных переиспользуемых пакетов, если репозиторий действительно предоставляет такой Go API. Engine сейчас организован как приложение, поэтому пустой `pkg/` не добавляем.

Зависимости направлены внутрь: доменные и прикладные пакеты не зависят от конкретных драйверов. Контракты размещаются рядом с кодом, который ими владеет; реализации — в adapter-пакетах на внешнем краю `internal/`.

## Разработка

Все команды проекта определены в `mise.toml`; отдельный Taskfile не используется.

- `mise install` — установка Go, Node.js и Go CLI-инструментов, зафиксированных в mise-конфигурации и lock-файлах.
- `mise run specs:prepare` — обновление TypeSpec/npm lock-файлов, установка npm-зависимостей и генерация OpenAPI.
- `mise run specs:gate` — проверка TypeSpec, lock-файлов и сгенерированного OpenAPI.
- `mise run sources:prepare` — генерация исходников и форматирование Go-кода.
- `mise run sources:gate` — проверка воспроизводимости генерации, форматирования, lint и static analysis.
- `mise run binaries:lock` — обновление Go module manifests и mise lock-файла.
- `mise run binaries:check` — проверка актуальности Go module manifests и mise lock-файла.
- `mise run binaries:build` — сборка бинарника Engine.
- `mise run binaries:unit-test` / `mise run binaries:integration-test` — unit- и интеграционные тесты.
- `mise run tests:e2e` — сквозные тесты.
- `mise run binaries:gate` — полный CI gate: Specs, Sources, lock checks, build и тесты.
- `mise run engine:process` — запуск приложения.
- `mise run db:process` — запуск PostgreSQL через Docker Compose.
- `mise run hooks:install` — настройка versioned Git hooks.

Версии Go, Node.js и Go CLI-инструментов фиксируются в `mise.toml`, `mise.lock` и stage-specific mise lock-файлах; зависимости приложения — в `go.mod` / `go.sum`; зависимости TypeSpec — в `tsp/package.json` / `tsp/package-lock.json`.
