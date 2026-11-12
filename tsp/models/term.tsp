import "./type.tsp";
import "./any/type.tsp";

namespace OrgLang.Engine.Term;

model VarSpec {
  chnl_ph: string;
  type_qn: string;
}

model VarRec {
  type_ref: Any.Type.Ref;
  chnl_ph: string;
  exp_vk: int64;
}
