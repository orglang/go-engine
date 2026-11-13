import "../../any/type.tsp";

namespace OrgLang.Pool.Type;

/** Pool type expression specification. Exactly one variant is serialized inline with its kind discriminator. */
@discriminated(#{ envelope: "none" })
union ExprSpec {
  one: OneSpec,
  link: LinkSpec,
  plus: PlusSpec,
  with: WithSpec,
  tensor: TensorSpec,
  lolli: LolliSpec,
  up: UpSpec,
  down: DownSpec,
}

model OneSpec {}
model LinkSpec { link: OrgLang.Any.Type.LinkSpec; }
model PlusSpec { proc_qns: string[]; cont_expr: ExprSpec; }
model WithSpec { proc_qns: string[]; cont_expr: ExprSpec; }
model TensorSpec { val_expr: ExprSpec; cont_expr: ExprSpec; }
model LolliSpec { val_expr: ExprSpec; cont_expr: ExprSpec; }
model UpSpec { cont_expr: ExprSpec; }
model DownSpec { cont_expr: ExprSpec; }

/** Additive fragment applied to an existing pool type expression. */
model DefPatch { extension: ExprSpec; }
