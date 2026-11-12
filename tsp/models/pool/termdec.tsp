import "../term.tsp";
import "../any/term.tsp";
namespace OrgLang.Pool;

model DecSpec {
  desc_qn: string;
  liab_var: Term.VarSpec;
  asset_vars: Term.VarSpec[];
}
