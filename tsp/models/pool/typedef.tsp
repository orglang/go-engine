import "../type.tsp";
import "../any/type.tsp";
import "./typeexp.tsp";
namespace OrgLang.Pool;

model DefSpec {
  qn: string;
  exp: TypeExprSpec;
}

model DefSnap {
  ref: Any.TypeRef;
  exp: TypeExprSpec;
}
