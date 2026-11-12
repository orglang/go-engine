import "../any/comp.tsp";

namespace OrgLang.Pool;

model VarSpec {
  conn_ph: string;
  term_qn: string;
}

model ExecSpec {
  term_qn: string;
  liab_var: VarSpec;
  asset_vars: VarSpec[];
}

model CompSnap {
  comp_ref: Any.CompRef;
}
