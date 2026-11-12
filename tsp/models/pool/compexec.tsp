import "../common.tsp";
namespace OrgLang.Engine;

model PoolExecSpec {
  term_qn: string;
  liab_var: ComputationVarSpec;
  asset_vars: ComputationVarSpec[];
}

model ComputationSnap {
  comp_ref: ComputationRef;
}
