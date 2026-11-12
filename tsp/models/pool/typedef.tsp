import "../common.tsp";
import "./typeexp.tsp";
namespace OrgLang.Engine.Pool;

model TypeSpec {
  qn: string;
  exp: TypeExprSpec;
}

model TypeSnap {
  ref: Common.TypeRef;
  exp: TypeExprSpec;
}
