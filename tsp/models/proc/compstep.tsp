import "../comp.tsp";
import "./termexp.tsp";

namespace OrgLang.Engine.Proc.Comp;

model StepSpec {
  comp_ref: OrgLang.Engine.Comp.Ref;
  proc_exp: OrgLang.Engine.Proc.TermExprSpec;
}
