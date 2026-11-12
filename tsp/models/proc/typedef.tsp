import "../type.tsp";
import "../any/type.tsp";
import "./typeexp.tsp";
namespace OrgLang.Proc;

model DefSpec {
  type_qn: string;
  type_exp: TypeExprSpec;
}

model DefSnap {
  ref: Any.TypeRef;
  spec: DefSpec;
}
