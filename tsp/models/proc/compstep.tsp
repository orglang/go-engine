import "../any/comp.tsp";
import "./termexp.tsp";

namespace OrgLang.Proc;

model StepSpec {
  comp_ref: Any.CompRef;
  proc_exp: Proc.TermExprSpec;
}
