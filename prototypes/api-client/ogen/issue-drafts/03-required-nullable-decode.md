## Description

ogen generates a decoder that rejects an explicit JSON `null` for an object property that is both required and nullable in an OpenAPI 3.1 schema.

Requiredness and nullability describe different aspects of a property:

- **Required** means that the property must be present in the JSON object.
- **Nullable** means that the property's value may be JSON `null`.

A required nullable property must therefore be present, but its value may be `null`.

## OpenAPI specification

The following OpenAPI 3.1 document reproduces the problem:

```yaml
openapi: 3.1.0
info:
  title: Required nullable recursive property
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
                $ref: '#/components/schemas/B'
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
          anyOf:
            - $ref: '#/components/schemas/A'
            - type: 'null'
```

The property `B.a` is required and nullable. Its nullable alternative is expressed using the OpenAPI 3.1 `anyOf` representation.

## Steps to reproduce

Save the document as `required-nullable.yaml`, generate a package, and decode the following JSON values into the generated `B` type:

```sh
go run github.com/ogen-go/ogen/cmd/ogen@v1.25.0 \
  --target /tmp/ogen-required-nullable \
  --package probe \
  --clean required-nullable.yaml
```

Use these inputs:

```json
{}
{"a":null}
{"a":{"b":{"a":null}}}
```

## Expected behavior

- `{}` is rejected because the required property `a` is absent.
- `{"a":null}` is accepted because the property is present and nullable.
- `{"a":{"b":{"a":null}}}` is accepted and decoded recursively.

## Actual behavior

The generated decoder rejects both inputs containing explicit `null`. It attempts to decode the value as an object and returns an error similar to:

```text
decode B: callback: decode field "a": decode A:
"{" expected: unexpected byte 110 'n'
```

The input with the missing property is rejected, as expected.

## Suspected cause

The generated decoder does not handle the null alternative of a nullable union when decoding the required property. It dispatches to the object decoder and rejects JSON `null`.

## Proposed change

When decoding a required nullable property, accept an explicit JSON `null` as a valid value while still marking the property as present for the required-field check. Keep rejecting the property when it is absent, and continue rejecting `null` for required non-nullable properties.

## Related issue

[ogen-go/ogen#1562](https://github.com/ogen-go/ogen/issues/1562) reports a related `null` decoding failure in the generated `OptString.Decode()` method. The case here is distinct: the property is required, its nullability is expressed as an OpenAPI 3.1 `anyOf` null alternative, and it participates in a recursive object schema. The required-property presence check must still reject an omitted property even though an explicit `null` is valid.
