import "../term.tsp";
import "../any/term.tsp";
namespace OrgLang.Proc;

model DecSpec {
  term_qn: string;
  liab_var: Term.VarSpec;
  asset_vars: Term.VarSpec[];
}

model DecSnap {
  term_ref: Any.TermRef;
  liab_var: Term.VarRec;
  asset_vars: Term.VarRec[];
}
