import "../any/comp.tsp";
import "./term/expr.tsp";

namespace OrgLang.Proc;

model ExecSnap {
  root_ref: Any.RootRef;
}

model StepSpec {
  proc_expr: Term.ExprSpec;
}
