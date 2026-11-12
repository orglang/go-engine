import "../../any/conn.tsp";
import "../../any/term.tsp";

namespace OrgLang.Pool;

model DecSpec {
  desc_qn: string;
  liab_var: Any.Conn.VarSpec;
  asset_vars: Any.Conn.VarSpec[];
}
