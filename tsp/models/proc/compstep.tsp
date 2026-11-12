import "../common.tsp";
import "./termexp.tsp";
namespace OrgLang.Engine.Proc;

model StepSpec {
  comp_ref: Common.ComputationRef;
  proc_exp: TermExprSpec;
}
