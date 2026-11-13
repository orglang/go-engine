import "../../any/conn/var.tsp";
import "../../any/term.tsp";

namespace OrgLang.Proc;

model DecSpec {
  term_qn: string;
  liab_var: Any.Conn.VarSpec;
  asset_vars: Any.Conn.VarSpec[];
}
