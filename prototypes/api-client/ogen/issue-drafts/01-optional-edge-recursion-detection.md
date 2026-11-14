# Bug: optional fields are incorrectly considered when detecting recursive Go value types

## Summary

`ogen` v1.25.0 may report an infinite-recursion error for a required field even when the apparent cycle is broken by an optional field.

Optional object properties are represented as pointers in generated Go code. A pointer edge breaks Go value recursion, so recursion detection must not traverse through an optional property when deciding whether an enclosing required field needs pointer boxing.

## Minimal OpenAPI 3.1 reproduction

Save as `optional-edge.yaml`:

```yaml
openapi: 3.1.0
info:
  title: Optional edge recursion
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

The graph is `A --required--> B --optional--> C --required--> B`. The optional `B.c` property is a pointer boundary, so `A.b` does not create an infinite Go value-recursion chain.

## Reproduction

With ogen v1.25.0, run:

```sh
go run github.com/ogen-go/ogen/cmd/ogen@v1.25.0 \
  --target /tmp/ogen-optional-edge \
  --package probe \
  --clean optional-edge.yaml
```

Expected: generation succeeds.

Actual: generation reports an error similar to:

```text
infinite recursion: A.b is required
```

## Proposed fix

In `gen/ir/recursion.go`, skip optional struct fields while walking the target type's fields:

```go
case KindStruct:
    return slices.ContainsFunc(target.Fields, func(f *Field) (r bool) {
        // Optional fields are represented through pointers and break
        // value-recursion.
        if f.Spec != nil && !f.Spec.Required {
            return false
        }
        return t.recursive(f.Type, path)
    })
```

This is the first hunk in the prototype patch:

[0001-ignore-optional-recursion-edges.patch](../patches/0001-ignore-optional-recursion-edges.patch)

## Scope

This is distinct from required-to-required recursion. It concerns a false positive caused by traversing an optional edge. The separate required-recursion case is documented in [02-required-recursive-fields.md](02-required-recursive-fields.md).
