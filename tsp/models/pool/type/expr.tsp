import "../../type.tsp";
import "../../any/type.tsp";

namespace OrgLang.Pool.Type;

/** Pool type expression specification. */
model ExprSpec {
  kind: "one" | "link" | "plus" | "with" | "tensor" | "lolli" | "up" | "down";
  link?: Type.LinkSpec;
  plus?: LaborSpec;
  with?: LaborSpec;
  tensor?: ResourceSpec;
  lolli?: ResourceSpec;
  up?: ShiftSpec;
  down?: ShiftSpec;
}

model LaborSpec { proc_qns: string[]; cont_expr: ExprSpec; }
model ResourceSpec { val_expr: ExprSpec; cont_expr: ExprSpec; }
model ShiftSpec { cont_expr: ExprSpec; }
