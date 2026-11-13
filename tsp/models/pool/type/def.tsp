import "../../any/type.tsp";
import "./expr.tsp";

namespace OrgLang.Pool.Type;

model DefSpec {
  qn: string;
  expr: ExprSpec;
}

model DefSnap {
  ref: Any.TypeRef;
  expr: ExprSpec;
}
