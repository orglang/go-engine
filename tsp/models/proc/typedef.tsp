import "../type.tsp";
import "../any/type.tsp";
import "./type/expr.tsp";

namespace OrgLang.Proc;

model DefSpec {
  type_qn: string;
  type_expr: Type.ExprSpec;
}

model DefSnap {
  ref: Any.TypeRef;
  spec: DefSpec;
}
