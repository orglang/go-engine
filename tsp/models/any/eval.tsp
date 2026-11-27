import "./eval/var.tsp";

namespace OrgLang.Any;

model EvalSpec {
  term_qn: string;
  liab_var: Eval.VarSpec;
  asset_vars: Eval.VarSpec[];
}
