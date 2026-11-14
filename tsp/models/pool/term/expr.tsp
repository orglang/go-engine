import "../../any/term.tsp";
import "../../any/comp.tsp";

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

// TODO(ogen-go/ogen#1790, ogen-go/ogen#1792): restore requiredness once ogen supports recursive optional edges and required-field pointer indirection.
model HireSpec { conn_ph: string; proc_qn: string; cont_expr?: ExprSpec; }
model FireSpec { conn_ph: string; }
// TODO(ogen-go/ogen#1790, ogen-go/ogen#1792): restore requiredness once ogen supports recursive optional edges and required-field pointer indirection.
model ApplySpec { conn_ph: string; proc_qn: string; cont_expr?: ExprSpec; }
model QuitSpec { conn_ph: string; proc_qn: string; }
// TODO(ogen-go/ogen#1790, ogen-go/ogen#1792): restore requiredness once ogen supports recursive optional edges and required-field pointer indirection.
model AcquireSpec { conn_ph: string; cont_expr?: ExprSpec; }
model ReleaseSpec { conn_ph: string; }
// TODO(ogen-go/ogen#1790, ogen-go/ogen#1792): restore requiredness once ogen supports recursive optional edges and required-field pointer indirection.
model AcceptSpec { conn_ph: string; cont_expr?: ExprSpec; }
model DetachSpec { conn_ph: string; }
model SpawnSpec { proc_dr: Any.TermRef; proc_irs: Any.CompRef[]; }
