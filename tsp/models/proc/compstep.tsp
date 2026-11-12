import "../any/comp.tsp";
import "./termexp.tsp";

namespace OrgLang.Engine.Proc;

model StepSpec {
  comp_ref: OrgLang.Engine.Any.CompRef;
  proc_exp: OrgLang.Engine.Proc.TermExprSpec;
}
