import "../../any/term.tsp";
import "../../any/comp.tsp";

namespace OrgLang.Pool.Term;

/** Pool term expression specification. */
model ExprSpec {
  kind: "hire" | "fire" | "apply" | "quit" | "acquire" | "release" | "accept" | "detach" | "spawn";
  hire?: HireSpec;
  fire?: FireSpec;
  apply?: ApplySpec;
  quit?: QuitSpec;
  acquire?: AcquireSpec;
  release?: ReleaseSpec;
  accept?: AcceptSpec;
  detach?: DetachSpec;
  spawn?: SpawnSpec;
}

model HireSpec { conn_ph: string; proc_qn: string; cont_expr: ExprSpec; }
model FireSpec { conn_ph: string; }
model ApplySpec { conn_ph: string; proc_qn: string; cont_expr: ExprSpec; }
model QuitSpec { conn_ph: string; proc_qn: string; }
model AcquireSpec { conn_ph: string; cont_expr: ExprSpec; }
model ReleaseSpec { conn_ph: string; }
model AcceptSpec { conn_ph: string; cont_expr: ExprSpec; }
model DetachSpec { conn_ph: string; }
model SpawnSpec { proc_dr: Any.TermRef; proc_irs: Any.CompRef[]; }
