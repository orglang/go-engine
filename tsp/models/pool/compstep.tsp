import "../common.tsp";
import "./termexp.tsp";
namespace OrgLang.Engine;

model PoolStepSpec {
  comp_ref: ComputationRef;
  pool_exp: PoolTermExprSpec;
}
