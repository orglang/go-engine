import "../type.tsp";
import "../term.tsp";
import "../comp.tsp";
import "./termexp.tsp";
namespace OrgLang.Engine.Proc;

model StepSpec {
  comp_ref: Comp.Ref;
  proc_exp: TermExprSpec;
}
