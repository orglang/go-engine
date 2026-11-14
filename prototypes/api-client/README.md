# Temporary Go client generator comparison

This directory is a disposable prototype for issue
[orglang/go-workspace#114](https://github.com/orglang/go-workspace/issues/114).

Both generators consume the same committed OpenAPI 3.1 document:

- `ogen v1.25.0`: generated transport and models in `ogen/`.
- `oapi-codegen v2.8.2`: generated models and HTTP client in `oapi-codegen/`.

## Domain-oriented ogen packages

The raw ogen package retains generator-produced names and is treated as an
implementation detail by consumers. Small domain packages provide a shorter,
stable surface:

```text
prototypes/api-client/
  ogen/       # generated transport, serialization and models
  pool/       # pool-facing aliases, constructors and methods
  process/    # process-facing aliases, constructors and methods
  program/    # program submission
```

The domain packages accept the shared generated client, so there is one HTTP
client/configuration and no duplicated transport implementation.

### Example: pool type

```go
package main

import (
    "context"
    "log"
    "net/http"
    "time"

    ogenclient "orglang/go-engine/prototypes/api-client/ogen"
    "orglang/go-engine/prototypes/api-client/pool"
)

func main() {
    httpClient := &http.Client{Timeout: 10 * time.Second}
    transport, err := ogenclient.NewClient(
        "http://localhost:8080",
        ogenclient.WithClient(httpClient),
    )
    if err != nil {
        log.Fatal(err)
    }

    pools := pool.New(transport)
    ctx := context.Background()

    created, err := pools.CreateType(ctx, &pool.TypeDefSpec{
        Qn:   "example.unit",
        Expr: ptr(pool.One()),
    })
    if err != nil {
        log.Fatal(err)
    }

    snapshot, err := pools.GetType(ctx, created.TypeId)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("loaded type revision: %d", snapshot.Ref.TypeRn)
}

func ptr[T any](v T) *T { return &v }
```

The domain API reads as `pool.CreateType`, `pool.GetType` and
`pool.One()`; it hides names such as
`EngineApiGetPoolTypeDefinition` and
`NewPoolTypeOneSpecPoolTypeExprSpec`.

### Example: process and program

```go
transport, err := ogenclient.NewClient("http://localhost:8080")
if err != nil {
    return err
}

processes := process.New(transport)
_, err = processes.CreateType(ctx, &process.TypeDefSpec{
    Qn: "example.process",
    Expr: ptr(process.One()),
})
if err != nil {
    return err
}

programs := program.New(transport)
if err := programs.Create(ctx, "process example.process = ..."); err != nil {
    return err
}
```

These examples illustrate the generated API shape; they do not call a running
engine. The facade currently focuses on the domain entry points and common
operations. Additional methods can be exposed without changing the generated
transport package.

## oapi-codegen example

`oapi-codegen` continues to expose the operation-oriented generated API:

```go
client, err := oapiclient.NewClientWithResponses(
    "http://localhost:8080",
    oapiclient.WithHTTPClient(httpClient),
)
if err != nil {
    return err
}

var expr oapiclient.PoolTypeExprSpec
if err := expr.FromPoolTypeOneSpec(oapiclient.PoolTypeOneSpec{}); err != nil {
    return err
}

created, err := client.EngineApiCreatePoolTypeDefinitionWithResponse(
    ctx,
    oapiclient.PoolTypeDefSpec{Qn: "example.unit", Expr: expr},
)
if err != nil {
    return err
}
if created.StatusCode() != http.StatusCreated || created.JSON201 == nil {
    return fmt.Errorf("create pool type: HTTP %s", created.Status())
}
```

## Comparison criteria

1. Complete contract generation and unsupported-schema behavior.
2. API ergonomics, optional/nullable values and discriminated unions.
3. Runtime dependencies and generated code size.
4. Customization, response handling and maintainability.
5. Reproducibility and generator version pinning.

No generated package here is intended for production use. All files in this
directory and the temporary task/workflow wiring will be removed after the
comparison is recorded.
