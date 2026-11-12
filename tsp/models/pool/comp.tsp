import "../any/comp.tsp";
import "../any/comp/var.tsp";
import "./term/expr.tsp";

namespace OrgLang.Pool;

model ExecSpec {
  term_qn: string;
  liab_var: Any.Comp.VarSpec;
  asset_vars: Any.Comp.VarSpec[];
}

model CompSnap {
  comp_ref: Any.CompRef;
}

model StepSpec {
  comp_ref: Any.CompRef;
  pool_expr: Term.ExprSpec;
}
