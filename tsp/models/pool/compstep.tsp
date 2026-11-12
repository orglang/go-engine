import "../any/comp.tsp";
import "./termexp.tsp";

namespace OrgLang.Engine.Pool;

model StepSpec {
  comp_ref: OrgLang.Engine.Any.CompRef;
  pool_exp: OrgLang.Engine.Pool.TermExprSpec;
}
