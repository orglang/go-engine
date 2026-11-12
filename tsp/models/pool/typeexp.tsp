import "../common.tsp";
namespace OrgLang.Engine;

/** Pool type expression specification (go-sdk/pool/typeexp). */
model PoolTypeExprSpec {
  kind: "one" | "link" | "plus" | "with" | "tensor" | "lolli" | "up" | "down";
  link?: TypeLinkSpec;
  plus?: PoolTypeLaborSpec;
  with?: PoolTypeLaborSpec;
  tensor?: PoolTypeResourceSpec;
  lolli?: PoolTypeResourceSpec;
  up?: PoolTypeShiftSpec;
  down?: PoolTypeShiftSpec;
}

model PoolTypeLaborSpec { proc_qns: string[]; cont_exp: PoolTypeExprSpec; }
model PoolTypeResourceSpec { val_exp: PoolTypeExprSpec; cont_exp: PoolTypeExprSpec; }
model PoolTypeShiftSpec { cont_exp: PoolTypeExprSpec; }
