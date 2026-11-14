## Description

ogen rejects recursive object schemas when the recursive edge is a required property:

```text
infinite recursion: ARecursive.b is required
```

A required property does not have to be represented as an embedded Go value. The generator can represent a recursive field through a pointer (or another indirection) while preserving the property's required status in the JSON contract.

This distinction matters because Go value types must have a finite size, whereas JSON objects can be recursively nested to arbitrary finite depth. A pointer provides the indirection needed for a finite Go representation; it does not make the JSON property optional.

## OpenAPI specification

The following OpenAPI 3.1 document is a minimal reproduction with a terminating alternative, so it admits valid finite JSON values:

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
      oneOf:
        - $ref: '#/components/schemas/ARecursive'
        - $ref: '#/components/schemas/ALeaf'
    ARecursive:
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
    ALeaf:
      type: object
      required: [kind]
      properties:
        kind:
          type: string
          enum: [leaf]
```

The recursive path is:

```text
ARecursive --required--> B --required--> A --oneOf--> ARecursive
```

The `ALeaf` alternative provides a finite termination case.

## Steps to reproduce

Save the document as `required-recursion.yaml` and run:

```sh
go run github.com/ogen-go/ogen/cmd/ogen@v1.25.0 \
  --target /tmp/ogen-required-recursion \
  --package probe \
  --clean required-recursion.yaml
```

## Expected behavior

Code generation succeeds by introducing indirection for a recursive edge in the generated Go type representation. The exact placement and shape of the pointer are implementation details.

The JSON contract remains unchanged:

- An object missing a required property is rejected.
- A present non-null recursive value is accepted.
- Explicit `null` is rejected for non-nullable properties.

For example, the decoder should accept both of these values:

```json
{"b":{"a":{"kind":"leaf"}}}
```

```json
{"b":{"a":{"b":{"a":{"kind":"leaf"}}}}}
```

It should reject an `ARecursive` value with its required `b` property omitted, for example:

```json
{}
```

## Actual behavior

ogen v1.25.0 fails during code generation with an infinite-recursion error for the required recursive property, rather than generating a finite Go representation.

## Suspected cause

The recursion check treats a required recursive property as an inlined Go value and aborts when it finds a cycle. Requiredness in OpenAPI describes whether a JSON object must contain a property; it does not require the generated Go field to be stored inline.

The generator appears to conflate two separate concerns:

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

## Compatibility considerations

This change should not alter the OpenAPI schema, its `required` arrays, or the accepted JSON values. It changes only the generated Go representation needed to make recursive types finite.

In Go, changing a field from a value to a pointer changes the in-memory representation and may affect source compatibility for consumers of generated types. That trade-off is separate from the wire contract: the pointer is an implementation detail of the generated model, while requiredness remains part of the API schema and decoder validation.
