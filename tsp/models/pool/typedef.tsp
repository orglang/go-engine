import "../type.tsp"
import "../any/type.tsp";;
import "./typeexp.tsp";
namespace OrgLang.Engine.Pool;

model TypeSpec {
  qn: string;
  exp: TypeExprSpec;
}

model TypeSnap {
  ref: Any.Type.Ref;
  exp: TypeExprSpec;
}
