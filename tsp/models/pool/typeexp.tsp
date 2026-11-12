import "../type.tsp";
namespace OrgLang.Engine.Pool;

/** Pool type expression specification (go-sdk/pool/typeexp). */
model TypeExprSpec {
  kind: "one" | "link" | "plus" | "with" | "tensor" | "lolli" | "up" | "down";
  link?: Type.LinkSpec;
  plus?: TypeLaborSpec;
  with?: TypeLaborSpec;
  tensor?: TypeResourceSpec;
  lolli?: TypeResourceSpec;
  up?: TypeShiftSpec;
  down?: TypeShiftSpec;
}

model TypeLaborSpec { proc_qns: string[]; cont_exp: TypeExprSpec; }
model TypeResourceSpec { val_exp: TypeExprSpec; cont_exp: TypeExprSpec; }
model TypeShiftSpec { cont_exp: TypeExprSpec; }
