import "../term.tsp";
import "../any/term.tsp";
namespace OrgLang.Engine.Pool;

model DeclSpec {
  desc_qn: string;
  liab_var: Term.VarSpec;
  asset_vars: Term.VarSpec[];
}
