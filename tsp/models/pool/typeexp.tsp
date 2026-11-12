import "../type.tsp";
import "../any/type.tsp";
namespace OrgLang.Pool;

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

model TypeLaborSpec { proc_qns: string[]; cont_expr: TypeExprSpec; }
model TypeResourceSpec { val_expr: TypeExprSpec; cont_expr: TypeExprSpec; }
model TypeShiftSpec { cont_expr: TypeExprSpec; }
