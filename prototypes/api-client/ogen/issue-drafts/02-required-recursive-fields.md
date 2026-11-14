# Bug: required recursive fields cannot be represented by generated Go types

## Summary

`ogen` v1.25.0 rejects schemas containing required recursive object properties with an `infinite recursion` error.

A required recursive property still needs a pointer in the generated Go type to make the Go type finite. The pointer must not make the property optional in JSON: its nil state must remain invalid.

## Minimal OpenAPI 3.1 reproduction

Save as `required-cycle.yaml`:

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
        default:
          description: Error
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
      required: [a]
      properties:
        a:
          $ref: '#/components/schemas/A'
```

Both edges are required. The generated Go representation needs a pointer on at least one recursive edge; otherwise the Go type has infinite size.

## Reproduction

With ogen v1.25.0, run:

```sh
go run github.com/ogen-go/ogen/cmd/ogen@v1.25.0 \
  --target /tmp/ogen-required-cycle \
  --package probe \
  --clean required-cycle.yaml
```

Expected: generation succeeds and the recursive Go field is pointer-backed while remaining required during JSON decoding/validation.

Actual: generation fails with an error similar to:

```text
infinite recursion: A.b is required
```

## Proposed fix

In `gen/generics.go`, when a recursive field is required, use a pointer with `ir.NilInvalid` instead of returning an error:

```go
default:
    // Required recursive fields need a pointer to break Go value
    // recursion, but nil remains invalid so JSON requiredness is kept.
    return ir.Pointer(t, ir.NilInvalid), nil
```

This keeps the generated Go type finite without changing the schema's requiredness semantics.

The corresponding second hunk is in:

[0001-ignore-optional-recursion-edges.patch](../patches/0001-ignore-optional-recursion-edges.patch)

## Scope

This is distinct from false-positive recursion through optional properties. The optional-edge case is documented in [01-optional-edge-recursion-detection.md](01-optional-edge-recursion-detection.md).
