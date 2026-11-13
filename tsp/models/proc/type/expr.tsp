import "../../any/type.tsp";

namespace OrgLang.Proc.Type;

/** Process type expression specification. Exactly one variant is serialized inline with its kind discriminator. */
@discriminated(#{ envelope: "none" })
union ExprSpec {
  one: OneSpec,
  link: LinkSpec,
  tensor: TensorSpec,
  lolli: LolliSpec,
  plus: PlusSpec,
  with: WithSpec,
  up: UpSpec,
  down: DownSpec,
}

model OneSpec {}
model LinkSpec { link: OrgLang.Any.Type.LinkSpec; }
model TensorSpec { val_expr: ExprSpec; cont_expr: ExprSpec; }
model LolliSpec { val_expr: ExprSpec; cont_expr: ExprSpec; }
model PlusSpec { choices: ChoiceSpec[]; }
model WithSpec { choices: ChoiceSpec[]; }
model ChoiceSpec { lab_qn: string; cont_expr: ExprSpec; }
model UpSpec { cont_expr: ExprSpec; }
model DownSpec { cont_expr: ExprSpec; }

/** Additive fragment applied to an existing process type expression. */
model DefPatch { extension: ExprSpec; }
