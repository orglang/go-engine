import "../type.tsp";
import "../any/type.tsp";
import "./type/expr.tsp";

namespace OrgLang.Pool;

model DefSpec {
  qn: string;
  expr: Type.ExprSpec;
}

model DefSnap {
  ref: Any.TypeRef;
  exp: Type.ExprSpec;
}
