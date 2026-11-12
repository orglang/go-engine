import "../../any/type.tsp";

namespace OrgLang.Proc.Type;

/** Process type expression specification. */
model ExprSpec {
  kind: "one" | "link" | "tensor" | "lolli" | "plus" | "with" | "up" | "down";
  link?: Any.Type.LinkSpec;
  tensor?: ProductSpec;
  lolli?: ProductSpec;
  plus?: SumSpec;
  with?: SumSpec;
  up?: ShiftSpec;
  down?: ShiftSpec;
}

model ProductSpec { val_exprr: ExprSpec; cont_exprr: ExprSpec; }
model SumSpec { choices: ChoiceSpec[]; }
model ChoiceSpec { lab_qn: string; cont_exprr: ExprSpec; }
model ShiftSpec { cont_exprr: ExprSpec; }
