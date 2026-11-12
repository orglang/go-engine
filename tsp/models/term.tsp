import "./type.tsp";
import "./any/type.tsp";

namespace OrgLang.Term;

model VarSpec {
  chnl_ph: string;
  type_qn: string;
}

model VarRec {
  type_ref: Any.TypeRef;
  chnl_ph: string;
  exp_ck: int64;
}
