import "./type.tsp";

namespace OrgLang.Engine.Term;

model Ref {
  term_id: string;
  term_rn: int64;
}

model VarSpec {
  chnl_ph: string;
  type_qn: string;
}

model VarRec {
  type_ref: Type.Ref;
  chnl_ph: string;
  exp_vk: int64;
}
