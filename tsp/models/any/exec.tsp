import "./exec/var.tsp";

namespace OrgLang.Any;

model ExecSpec {
  term_qn: string;
  liab_var: Exec.VarSpec;
  asset_vars: Exec.VarSpec[];
}
