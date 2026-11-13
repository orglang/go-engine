import "../any/comp.tsp";
import "./term/expr.tsp";

namespace OrgLang.Proc;

model CompSnap {
  comp_ref: Any.CompRef;
}

model StepSpec {
  proc_expr: Term.ExprSpec;
}
