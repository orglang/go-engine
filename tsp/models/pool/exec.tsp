import "../any/comp.tsp";
import "./term/expr.tsp";

namespace OrgLang.Pool;

model ExecSnap {
  root_ref: Any.RootRef;
}

model StepSpec {
  root_ref: Any.RootRef;
  pool_expr: Term.ExprSpec;
}
