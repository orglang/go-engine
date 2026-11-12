import "../comp.tsp";

namespace OrgLang.Engine.Pool.Comp;

model VarSpec {
  chnl_ph: string;
  term_qn: string;
}

model ExecSpec {
  term_qn: string;
  liab_var: VarSpec;
  asset_vars: VarSpec[];
}

model Snap {
  comp_ref: OrgLang.Engine.Comp.Ref;
}
