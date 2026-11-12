import "../type.tsp";
import "../any/type.tsp";
import "./typeexp.tsp";
namespace OrgLang.Proc;

model DefSpec {
  type_qn: string;
  type_expr: TypeExprSpec;
}

model DefSnap {
  ref: Any.TypeRef;
  spec: DefSpec;
}
