import "../term.tsp";
import "../any/term.tsp";
namespace OrgLang.Engine.Proc;

model DeclSpec {
  term_qn: string;
  liab_var: Term.VarSpec;
  asset_vars: Term.VarSpec[];
}

model DeclSnap {
  term_ref: Any.Term.Ref;
  liab_var: Term.VarRec;
  asset_vars: Term.VarRec[];
}
