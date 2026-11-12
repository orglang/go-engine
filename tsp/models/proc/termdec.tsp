import "../common.tsp";
namespace OrgLang.Engine;

model ProcDeclSpec {
  term_qn: string;
  liab_var: TermVarSpec;
  asset_vars: TermVarSpec[];
}

model ProcDeclSnap {
  term_ref: TermRef;
  liab_var: TermVarRec;
  asset_vars: TermVarRec[];
}
