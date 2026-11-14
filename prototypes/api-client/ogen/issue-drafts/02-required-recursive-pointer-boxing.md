## Description

ogen rejects recursive object schemas when the recursive edge is a required property:

```text
infinite recursion: A.b is required
```

A required property does not have to be represented as an embedded Go value. The generator can represent a recursive field through a pointer (or another indirection) while preserving the property's required status in the JSON contract.

This distinction matters because Go value types must have a finite size, whereas JSON objects can be recursively nested to arbitrary finite depth. A pointer provides the indirection needed for a finite Go representation; it does not make the JSON property optional.

## OpenAPI specification

The following OpenAPI 3.1 document is a minimal reproduction:

```yaml
openapi: 3.1.0
info:
  title: Required recursive properties
  version: 1.0.0
paths:
  /probe:
    post:
      operationId: probe
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/A'
      responses:
        '204':
          description: No content
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

Both properties are required and non-nullable. The recursive cycle is:

```text
A --required--> B --required--> A
```

## Steps to reproduce

Save the document as `required-recursion.yaml` and run:

```sh
go run github.com/ogen-go/ogen/cmd/ogen@v1.25.0 \
  --target /tmp/ogen-required-recursion \
  --package probe \
  --clean required-recursion.yaml
```

## Expected behavior

Code generation succeeds by introducing indirection for at least one edge in the recursive Go type representation. For example, the generated types may be equivalent to:

```go
type A struct {
    B *B
}

type B struct {
    A *A
}
```

The exact placement and shape of the pointer are an implementation detail. The JSON contract must remain unchanged:

- An object missing a required property is rejected.
- A present object value is accepted, including recursively nested values.
- Explicit `null` is rejected because these properties are non-nullable.

For example, the decoder should reject:

```json
{}
```

and accept a finite nested value such as:

```json
{"b":{"a":{"b":{"a":{}}}}}
```

The innermost object in this example is intentionally incomplete only if it is decoded as `A`; a fully valid finite value can instead terminate at a different schema branch. For this two-type reproduction, both properties are required, so every finite JSON value must eventually be rejected unless the schema provides a terminating alternative. The code-generation expectation is independent of whether this particular schema admits a finite valid instance.

## Actual behavior

ogen v1.25.0 fails during code generation with an infinite-recursion error for the required recursive property, rather than generating a finite Go representation.

## Suspected cause

The recursion check treats a required recursive property as an inlined Go value and aborts when it finds a cycle. Requiredness in OpenAPI describes whether a JSON object must contain a property; it does not require the generated Go field to be stored inline.

The generator appears to conflate these two concerns:

- **Wire-level requiredness:** whether the JSON object contains the property.
- **Go representation:** whether the generated field is a value or an indirection such as a pointer.

They should be handled independently.

## Proposed change

When a required struct field closes a value-recursion cycle, generate that field through pointer indirection instead of returning an infinite-recursion error.

The generated decoder must continue tracking whether the property appeared in the input. In particular, a nil Go pointer must not by itself be used to infer that the property was absent: required-field validation should be based on JSON field presence.

A regression test should verify both generated type finiteness and JSON behavior. It should cover:

1. Missing required recursive property is rejected.
2. Present non-null recursive property is decoded.
3. Explicit `null` is rejected for a non-nullable property.
4. Marshaling and decoding a valid finite instance preserves the recursive structure.

A schema with a terminating alternative (for example, a recursive `oneOf` with a non-recursive leaf variant) is suitable for the runtime round-trip test, while the two-object schema above isolates the code-generation failure.

## Compatibility considerations

This change should not alter the OpenAPI schema, its `required` arrays, or the accepted JSON values. It changes only the generated Go representation needed to make recursive types finite.

In Go, changing a field from a value to a pointer changes the in-memory representation and may affect source compatibility for consumers of generated types. That trade-off is separate from the wire contract: the pointer is an implementation detail of the generated model, while requiredness remains part of the API schema and decoder validation.
