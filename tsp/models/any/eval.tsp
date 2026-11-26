import "./eval/var.tsp";

namespace OrgLang.Any;

model EvalSpec {
  term_qn: string;
  liab_var: Exec.VarSpec;
  asset_vars: Exec.VarSpec[];
}
