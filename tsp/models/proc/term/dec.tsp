import "../../any/conn.tsp";
import "../../any/term.tsp";

namespace OrgLang.Proc;

model DecSpec {
  term_qn: string;
  liab_var: Any.Conn.VarSpec;
  asset_vars: Any.Conn.VarSpec[];
}

model DecSnap {
  term_ref: Any.TermRef;
  liab_var: Any.Conn.VarRec;
  asset_vars: Any.Conn.VarRec[];
}
