import "../any/comp.tsp";
import "./termexp.tsp";

namespace OrgLang.Pool;

model StepSpec {
  comp_ref: Any.CompRef;
  pool_exp: Pool.TermExprSpec;
}
