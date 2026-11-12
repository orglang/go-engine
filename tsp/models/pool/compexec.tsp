import "../common.tsp";
namespace OrgLang.Engine.Pool;

model ExecSpec {
  term_qn: string;
  liab_var: Common.CompVarSpec;
  asset_vars: Common.CompVarSpec[];
}

model ComputationSnap {
  comp_ref: Common.ComputationRef;
}
