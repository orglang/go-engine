import "../common.tsp";
namespace OrgLang.Engine.Proc;

model DeclSpec {
  term_qn: string;
  liab_var: Common.TermVarSpec;
  asset_vars: Common.TermVarSpec[];
}

model DeclSnap {
  term_ref: Common.TermRef;
  liab_var: Common.TermVarRec;
  asset_vars: Common.TermVarRec[];
}
