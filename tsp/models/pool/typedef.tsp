import "./typeexp.tsp";
namespace OrgLang.Engine;

model PoolTypeSpec {
  qn: string;
  exp: PoolTypeExprSpec;
}

model PoolTypeSnap {
  ref: TypeRef;
  exp: PoolTypeExprSpec;
}
