import "../common.tsp";
namespace OrgLang.Engine;

model PoolDeclSpec {
  desc_qn: string;
  liab_var: TermVarSpec;
  asset_vars: TermVarSpec[];
}
