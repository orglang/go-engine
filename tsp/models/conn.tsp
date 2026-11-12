import "./any/type.tsp";

namespace OrgLang.Conn;

model VarSpec {
  conn_ph: string;
  type_qn: string;
}

model VarRec {
  type_ref: Any.TypeRef;
  conn_ph: string;
  exp_ck: int64;
}
