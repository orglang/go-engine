import "../../any/term.tsp";
import "../../any/root.tsp";

namespace OrgLang.Pool.Term;

/** Pool term expression specification. Exactly one variant is serialized inline with its kind discriminator. */
@discriminated(#{ envelope: "none" })
union ExprSpec {
  hire: HireSpec,
  fire: FireSpec,
  apply: ApplySpec,
  quit: QuitSpec,
  acquire: AcquireSpec,
  release: ReleaseSpec,
  accept: AcceptSpec,
  detach: DetachSpec,
  spawn: SpawnSpec,
}

model HireSpec {
  comm_ph: string;
  proc_qn: string;
  cont_expr: ExprSpec;
}
model FireSpec {
  comm_ph: string;
}
model ApplySpec {
  comm_ph: string;
  proc_qn: string;
  cont_expr: ExprSpec;
}
model QuitSpec {
  comm_ph: string;
  proc_qn: string;
}
model AcquireSpec {
  comm_ph: string;
  cont_expr: ExprSpec;
}
model ReleaseSpec {
  comm_ph: string;
}
model AcceptSpec {
  comm_ph: string;
  cont_expr: ExprSpec;
}
model DetachSpec {
  comm_ph: string;
}
model SpawnSpec {
  proc_dr: Any.TermRef;
  proc_exec_refs: Any.RootRef[];
}
