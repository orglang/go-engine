import "../../any/type.tsp";
import "./expr.tsp";

namespace OrgLang.Proc.Type;

model DefSpec {
  type_qn: string;
  type_expr: ExprSpec;
}

model DefSnap {
  ref: Any.TypeRef;
  spec: DefSpec;
}
