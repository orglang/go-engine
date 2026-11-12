import "../comp.tsp";
import "./termexp.tsp";
namespace OrgLang.Engine.Pool;

model StepSpec {
  comp_ref: Comp.Ref;
  pool_exp: TermExprSpec;
}
