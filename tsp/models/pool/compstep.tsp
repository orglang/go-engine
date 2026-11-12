import "../common.tsp";
import "./termexp.tsp";
namespace OrgLang.Engine.Pool;

model StepSpec {
  comp_ref: Common.ComputationRef;
  pool_exp: TermExprSpec;
}
