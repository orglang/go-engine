import "../type.tsp";
import "../term.tsp";
import "../comp.tsp";
import "./typeexp.tsp";
namespace OrgLang.Engine.Pool;

model TypeSpec {
  qn: string;
  exp: TypeExprSpec;
}

model TypeSnap {
  ref: Type.Ref;
  exp: TypeExprSpec;
}
