import "./typeexp.tsp";
namespace OrgLang.Engine;

model ProcTypeSpec {
  type_qn: string;
  type_exp: ProcTypeExprSpec;
}

model ProcTypeSnap {
  ref: TypeRef;
  spec: ProcTypeSpec;
}
