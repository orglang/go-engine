import "@typespec/http";
import "@typespec/openapi3";

using TypeSpec.Http;

@service(#{ title: "Required nullable recursive TypeSpec probe" })
namespace RequiredNullableProbe;

model A {
  b: B;
}

model B {
  a: A | null;
}

@get
@route("/probe")
op probe(): A;
