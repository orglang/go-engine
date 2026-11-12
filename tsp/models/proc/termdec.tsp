import "../conn.tsp";
import "../any/term.tsp";

namespace OrgLang.Proc;

model DecSpec {
  term_qn: string;
  liab_var: Conn.VarSpec;
  asset_vars: Conn.VarSpec[];
}

model DecSnap {
  term_ref: Any.TermRef;
  liab_var: Conn.VarRec;
  asset_vars: Conn.VarRec[];
}
