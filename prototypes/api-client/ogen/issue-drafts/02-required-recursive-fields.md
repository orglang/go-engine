## Description

ogen fails to generate code for a valid schema when a required object property participates in a recursive cycle.

Go structs cannot contain themselves, directly or indirectly, by value. A recursive required property therefore needs indirection in the generated Go representation. That implementation detail must not change the JSON contract: the property remains required and must be rejected when absent.

## OpenAPI specification

The following OpenAPI 3.1 document reproduces the problem:

```yaml
openapi: 3.1.0
info:
  title: Required recursive fields
  version: 1.0.0
paths:
  /probe:
    get:
      operationId: probe
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/A'
components:
  schemas:
    A:
      type: object
      required: [b]
      properties:
        b:
          $ref: '#/components/schemas/B'
    B:
      type: object
      required: [c]
      properties:
        c:
          $ref: '#/components/schemas/C'
    C:
      type: object
      required: [b]
      properties:
        b:
          $ref: '#/components/schemas/B'
```

The recursive cycle is:

```text
B --required--> C --required--> B
```

Both edges are required; there is no optional property to break the value-recursion cycle.

## Steps to reproduce

Save the document as `required-cycle.yaml` and run:

```sh
go run github.com/ogen-go/ogen/cmd/ogen@v1.25.0 \
  --target /tmp/ogen-required-cycle \
  --package probe \
  --clean required-cycle.yaml
```

## Expected behavior

Code generation succeeds and emits finite Go types. The generated representation uses pointer indirection (or an equivalent mechanism) for at least one edge in the recursive cycle.

The JSON semantics remain unchanged: decoding an object that omits a required recursive property must fail.

## Actual behavior

Code generation fails with an infinite-recursion error similar to:

```text
infinite recursion: B.c is required
```

## Suspected cause

When the generator detects a required recursive field, it reports an error instead of selecting an indirect Go representation for that field.

The generator should break the Go value-recursion cycle using indirection while preserving the field's required semantics during decoding and validation.
