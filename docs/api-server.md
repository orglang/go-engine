# Echo HTTP server generation

The HTTP API contract is maintained in TypeSpec under `tsp/`. The generated OpenAPI document is committed at `openapi/openapi.yaml`. The Echo server boundary and request/response models are generated from that document into `internal/app/web/api/api.gen.go`.

## Generate

Run the full source-generation pipeline:

```sh
task sources:generate
```

This runs the Go source generators, including `api:generate`. To run only the API generator:

```sh
task api:generate
```

The generator version is pinned in `gear.mod`, and its configuration is in `internal/app/web/api/cfg.yaml`. Do not edit `api.gen.go` by hand; change the TypeSpec contract and regenerate the OpenAPI document and Go sources instead.

## Use the generated Echo boundary

Implement `api.ServerInterface` for the application handler and register it on an Echo instance:

```go
e := echo.New()
api.RegisterHandlers(e, handler)
```

The generated registration binds path and query parameters and performs the generated request decoding/validation before calling the handler. The integration test in `test/it/api_test.go` demonstrates route registration, parameter binding, JSON response handling, and rejection of an invalid integer query parameter.

## Check generated output

Run:

```sh
task sources:check
```

The check regenerates sources and fails if regeneration changes the committed `internal/app/web/api/api.gen.go`. It then runs formatting and lint checks. CI invokes `sources:check`, so API generation drift is checked with the rest of the source pipeline.
