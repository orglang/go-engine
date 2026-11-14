## Description

The recursion detector reports infinite recursion for a required property even when the recursive path is broken by an optional property.

Optional object properties are represented by pointer-like optional wrappers in generated Go types. Such a property breaks Go value recursion and must not be followed as a value edge when checking whether a type is recursively embedded by value.

## OpenAPI specification

The following OpenAPI 3.1 document reproduces the problem:

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

The relevant graph is:

```text
A --required--> B --optional--> C --required--> B
```

The optional `B.c` edge breaks the Go value-recursion cycle.

## Steps to reproduce

Save the document as `optional-edge.yaml` and run:

```sh
go run github.com/ogen-go/ogen/cmd/ogen@v1.25.0 \
  --target /tmp/ogen-optional-edge \
  --package probe \
  --clean optional-edge.yaml
```

## Expected behavior

Code generation succeeds. The generated Go types are finite because the optional property provides an indirection in the recursive path.

## Actual behavior

Code generation fails with an infinite-recursion error similar to:

```text
infinite recursion: A.b is required
```

## Suspected cause

The recursion detection for struct types traverses all fields without excluding optional fields. It therefore treats an optional pointer-like edge as a Go value edge and reports a cycle that cannot cause infinite value recursion.

## Proposed change

Exclude optional struct fields from the value-recursion graph. Their generated representation already provides indirection, so they cannot create infinite recursion through Go struct values. Continue traversing required fields so that cycles consisting entirely of value fields are still detected.

## Related issue

[ogen-go/ogen#1488](https://github.com/ogen-go/ogen/issues/1488) also concerns recursive schema generation involving `allOf`, but reports a different failure: the generated `OptError` type is missing. This report is specifically about the recursion detector treating an optional edge as a value edge and incorrectly rejecting a cycle that is already broken by indirection.
