# Temporary Go client generator comparison

This directory is a disposable prototype for issue
[orglang/go-workspace#114](https://github.com/orglang/go-workspace/issues/114).

Both generators consume the same committed OpenAPI 3.1 document:

- `ogen v1.25.0`: client-only generation, configured in `ogen.yml`.
- `oapi-codegen v2.8.2`: generated models and HTTP client.

Run the prototype through the repository's existing Generate artifacts
workflow (or locally with `mise run generate` on this branch).

## Usage examples

The snippets below use the generated packages in this directory. They show
creating a pool type definition whose expression is the `one` variant, then
fetching a type definition. They are examples of generated API shape, not calls
to a running engine.

### ogen

```go
package main

import (
    "context"
    "log"
    "net/http"
    "time"

    ogenclient "orglang/go-engine/prototypes/api-client/ogen"
)

func main() {
    httpClient := &http.Client{Timeout: 10 * time.Second}

    client, err := ogenclient.NewClient(
        "http://localhost:8080",
        ogenclient.WithClient(httpClient),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Build the discriminated union with the generated constructor.
    expr := ogenclient.NewPoolTypeOneSpecPoolTypeExprSpec(
        ogenclient.PoolTypeOneSpec{},
    )

    created, err := client.EngineApiCreatePoolTypeDefinition(
        context.Background(),
        &ogenclient.PoolTypeDefSpec{
            Qn:   "example.unit",
            Expr: &expr,
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("created type: %s (revision %d)", created.TypeId, created.TypeRn)

    // Path parameter and optional query parameter are grouped in Params.
    snapshot, err := client.EngineApiGetPoolTypeDefinition(
        context.Background(),
        ogenclient.EngineApiGetPoolTypeDefinitionParams{
            ID: created.TypeId,
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("loaded type revision: %d", snapshot.Ref.TypeRn)
}
```

Key characteristics:

- The generated client owns request construction, JSON encoding/decoding and
  response status handling.
- `oneOf` is represented as a tagged Go sum type. Use generated constructors
  and accessors rather than assembling raw JSON.
- The generated client supports request editors/options and OpenTelemetry
  instrumentation (depending on enabled generator features).

### oapi-codegen

```go
package main

import (
    "context"
    "log"
    "net/http"
    "time"

    oapiclient "orglang/go-engine/prototypes/api-client/oapi-codegen"
)

func main() {
    httpClient := &http.Client{Timeout: 10 * time.Second}

    client, err := oapiclient.NewClientWithResponses(
        "http://localhost:8080",
        oapiclient.WithHTTPClient(httpClient),
    )
    if err != nil {
        log.Fatal(err)
    }

    // The generated union type exposes From... helpers. The helper inserts
    // the discriminator into the serialized JSON representation.
    var expr oapiclient.PoolTypeExprSpec
    if err := expr.FromPoolTypeOneSpec(oapiclient.PoolTypeOneSpec{}); err != nil {
        log.Fatal(err)
    }

    created, err := client.EngineApiCreatePoolTypeDefinitionWithResponse(
        context.Background(),
        oapiclient.PoolTypeDefSpec{
            Qn:   "example.unit",
            Expr: expr,
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    if created.StatusCode() != http.StatusCreated {
        log.Fatalf("create type: HTTP %s: %s", created.Status(), string(created.Body))
    }
    if created.JSON201 == nil {
        log.Fatal("create type: missing JSON response")
    }
    typeID := created.JSON201.TypeId

    snapshot, err := client.EngineApiGetPoolTypeDefinitionWithResponse(
        context.Background(),
        typeID,
        nil, // omit rn to request the latest revision
    )
    if err != nil {
        log.Fatal(err)
    }
    if snapshot.StatusCode() != http.StatusOK {
        log.Fatalf("get type: HTTP %s: %s", snapshot.Status(), string(snapshot.Body))
    }
    if snapshot.JSON200 == nil {
        log.Fatal("get type: missing JSON response")
    }
    log.Printf("loaded type revision: %d", snapshot.JSON200.Ref.TypeRn)
}
```

Key characteristics:

- The generated client returns typed response wrappers; callers inspect status
  and the response variant (for example `JSON200`) explicitly.
- Discriminated unions are backed by raw JSON internally. The generated
  `From...` and `As...` helpers make construction and decoding possible, but
  callers work with a less idiomatic representation than ogen's sum type.
- The runtime dependency is `github.com/oapi-codegen/runtime`.

## Comparison criteria

1. Whether the complete contract generates without unsupported-schema errors.
2. Generated API ergonomics: operation names, request/response types, optional
   and nullable values, and discriminated expression unions.
3. Runtime dependencies and generated code size.
4. Error/response handling, customization points, and maintainability.
5. Reproducibility and ease of pinning the generator version.

No generated package here is intended for production use. All files in this
directory and the temporary task/workflow wiring will be removed after the
comparison is recorded.
