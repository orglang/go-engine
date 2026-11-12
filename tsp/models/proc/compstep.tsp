import "../common.tsp";
import "./termexp.tsp";
namespace OrgLang.Engine;

model ProcStepSpec {
  comp_ref: ComputationRef;
  proc_exp: ProcTermExprSpec;
}
