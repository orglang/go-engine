import "../type.tsp";
import "./typeexp.tsp";
namespace OrgLang.Engine.Proc;

model TypeSpec {
  type_qn: string;
  type_exp: TypeExprSpec;
}

model TypeSnap {
  ref: Type.Ref;
  spec: TypeSpec;
}
