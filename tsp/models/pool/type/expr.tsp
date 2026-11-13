import "../../any/type.tsp";

namespace OrgLang.Pool.Type;

/** Pool type expression specification. Exactly one variant is serialized inline with its kind discriminator. */
@discriminated(#{ envelope: "none" })
union ExprSpec {
  one: OneSpec,
  link: LinkSpec,
  plus: LaborSpec,
  with: LaborSpec,
  tensor: ResourceSpec,
  lolli: ResourceSpec,
  up: ShiftSpec,
  down: ShiftSpec,
}

model OneSpec {}
model LinkSpec { link: OrgLang.Any.Type.LinkSpec; }
model LaborSpec { proc_qns: string[]; cont_expr: ExprSpec; }
model ResourceSpec { val_expr: ExprSpec; cont_expr: ExprSpec; }
model ShiftSpec { cont_expr: ExprSpec; }

/** Additive fragment applied to an existing pool type expression. */
model DefPatch { extension: ExprSpec; }
