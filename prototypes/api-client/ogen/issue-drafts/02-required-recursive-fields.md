# Bug: ogen rejects required recursive object fields during generation

## Summary

ogen v1.25.0 rejects a schema when a required object property participates in a recursive cycle. The schema is valid, but the generator reports infinite recursion before it can emit Go code.

The generated Go representation must be finite. At least one edge in a required-to-required cycle therefore needs pointer indirection, while the JSON property must remain required: omitting it must still be rejected by decoding/validation.

## Minimal OpenAPI 3.1 reproduction

Save as required-cycle.yaml:

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

The recursive cycle is B --required--> C --required--> B. Both properties are required; there is no optional pointer edge to break the Go value cycle.

## Reproduction

With ogen v1.25.0, run:

```sh
go run github.com/ogen-go/ogen/cmd/ogen@v1.25.0 \
  --target /tmp/ogen-required-cycle \
  --package probe \
  --clean required-cycle.yaml
```


**Expected:** generation succeeds and emits finite Go types. The generated representation uses pointer indirection for a recursive edge, without treating the JSON property as optional.

**Actual:** generation fails with an error similar to:

```text
infinite recursion: B.c is required
```


This failure is reproducible with the unmodified v1.25.0 generator.

## Prototype fix

The prototype changes the required-recursive-field fallback in gen/generics.go from returning an error to returning a pointer type with ir.NilInvalid:

```go
default:
    // Required recursive fields need a pointer to break Go value
    // recursion, but nil remains invalid so JSON requiredness is kept.
    return ir.Pointer(t, ir.NilInvalid), nil
```


The intent is to make the Go type finite while retaining required-field validation. In the repository prototype, both the required-cycle fixture and the TypeSpec-generated required-recursion probe complete generation with the patch; unpatched ogen fails on the required recursive cycle.

This is a generator/type-representation fix proposal. Runtime behavior must still be tested separately: successful code generation alone does not prove that missing, null, and non-null JSON values are handled correctly.

## Nullable control experiment and separate finding

The repository also contains fixtures/required-nullable-cycle.yaml, plus a TypeSpec 1.18.0 probe that emits an OpenAPI 3.1 schema with a required nullable recursive property.

The patched and unpatched ogen v1.25.0 generators both successfully generate this nullable probe. However, the runtime test shows that **both generated clients reject an explicit null** for the required nullable property. Both also reject the nested object containing that null. The test separately confirms that a missing property is rejected.

Because the unpatched generator has the same runtime behavior, this nullable decoding problem is not evidence of a regression introduced by the required-recursion prototype. It is a separate upstream behavior/bug concerning the distinction between requiredness and nullability. It should be tracked independently rather than used to broaden this issue.

The nullable control is useful because it demonstrates that NilInvalid by itself does not establish correct nullable semantics. Do not claim that the proposed pointer fallback has fully preserved all JSON semantics until nullable decoding is addressed.

## Scope

This issue is specifically about generation failing on a required-to-required recursive object cycle and the need for finite Go type representation. It is distinct from:

- false-positive recursion detection through optional fields, documented in 01-optional-edge-recursion-detection.md;
- incorrect handling of explicit null for required nullable properties, observed in the nullable control experiment above and requiring a separate report.
