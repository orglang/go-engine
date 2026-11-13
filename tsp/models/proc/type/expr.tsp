import "../../any/type.tsp";

namespace OrgLang.Proc.Type;

/** Process type expression specification. Exactly one variant is serialized inline with its kind discriminator. */
@discriminated(#{ envelope: "none" })
union ExprSpec {
  one: OneSpec,
  link: LinkSpec,
  tensor: ProductSpec,
  lolli: ProductSpec,
  plus: SumSpec,
  with: SumSpec,
  up: ShiftSpec,
  down: ShiftSpec,
}

model OneSpec {}
model LinkSpec { link: OrgLang.Any.Type.LinkSpec; }
model ProductSpec { val_expr: ExprSpec; cont_expr: ExprSpec; }
model SumSpec { choices: ChoiceSpec[]; }
model ChoiceSpec { lab_qn: string; cont_expr: ExprSpec; }
model ShiftSpec { cont_expr: ExprSpec; }

/** Additive fragment applied to an existing process type expression. */
model DefPatch { extension: ExprSpec; }
