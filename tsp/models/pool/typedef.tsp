import "../type.tsp";
import "../any/type.tsp";
import "./typeexp.tsp";
namespace OrgLang.Pool;

model DefSpec {
  qn: string;
  expr: TypeExprSpec;
}

model DefSnap {
  ref: Any.TypeRef;
  exp: TypeExprSpec;
}
