import "../conn.tsp";
import "../any/term.tsp";

namespace OrgLang.Pool;

model DecSpec {
  desc_qn: string;
  liab_var: Conn.VarSpec;
  asset_vars: Conn.VarSpec[];
}
