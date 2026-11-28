import "../../any/comm/var.tsp";
import "../../any/term.tsp";

namespace OrgLang.Pool;

model DecSpec {
  desc_qn: string;
  liab_var: Any.Comm.VarSpec;
  asset_vars: Any.Comm.VarSpec[];
}

model DecSnap {
  ref: Any.TermRef;
  spec: DecSpec;
}
