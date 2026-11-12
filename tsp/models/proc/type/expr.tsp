import "../../any/type.tsp";

namespace OrgLang.Proc.Type;

/** Process type expression specification. */
model ExprSpec {
  kind: "one" | "link" | "tensor" | "lolli" | "plus" | "with" | "up" | "down";
  link?: OrgLang.Any.Type.LinkSpec;
  tensor?: ProductSpec;
  lolli?: ProductSpec;
  plus?: SumSpec;
  with?: SumSpec;
  up?: ShiftSpec;
  down?: ShiftSpec;
}

model ProductSpec { val_expr: ExprSpec; cont_expr: ExprSpec; }
model SumSpec { choices: ChoiceSpec[]; }
model ChoiceSpec { lab_qn: string; cont_expr: ExprSpec; }
model ShiftSpec { cont_expr: ExprSpec; }
