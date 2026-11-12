import "../comp.tsp";
import "./termexp.tsp";

namespace OrgLang.Engine.Pool.Comp;

model StepSpec {
  comp_ref: OrgLang.Engine.Comp.Ref;
  pool_exp: OrgLang.Engine.Pool.TermExprSpec;
}
