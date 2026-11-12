import "../common.tsp";
namespace OrgLang.Engine.Pool;

model DeclSpec {
  desc_qn: string;
  liab_var: Common.TermVarSpec;
  asset_vars: Common.TermVarSpec[];
}
