import "../any/comp.tsp";

namespace OrgLang.Engine.Pool;

model VarSpec {
  chnl_ph: string;
  term_qn: string;
}

model ExecSpec {
  term_qn: string;
  liab_var: VarSpec;
  asset_vars: VarSpec[];
}

model CompSnap {
  comp_ref: OrgLang.Engine.Any.Comp.Ref;
}
