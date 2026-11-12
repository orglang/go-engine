import "../any/comp.tsp";
import "./term/expr.tsp";

namespace OrgLang.Pool;

model StepSpec {
  comp_ref: Any.CompRef;
  pool_expr: Term.ExprSpec;
}
