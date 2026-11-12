import "../term.tsp";
import "../any/term.tsp";
import "../any/comp.tsp";
namespace OrgLang.Engine.Pool;

/** Pool term expression specification (go-sdk/pool/termexp). */
model TermExprSpec {
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

model HireSpec { comm_ph: string; proc_qn: string; cont_exp: TermExprSpec; }
model FireSpec { CommChnlPH: string; }
model ApplySpec { comm_ph: string; proc_qn: string; cont_exp: TermExprSpec; }
model QuitSpec { comm_ph: string; proc_qn: string; }
model AcquireSpec { comm_ph: string; cont_exp: TermExprSpec; }
model ReleaseSpec { comm_ph: string; }
model AcceptSpec { comm_ph: string; cont_exp: TermExprSpec; }
model DetachSpec { comm_ph: string; }
model SpawnSpec { proc_dr: Any.TermRef; proc_irs: Any.CompRef[]; }
