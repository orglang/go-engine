import "../type.tsp";

namespace OrgLang.Any.Conn;

model VarSpec {
  conn_ph: string;
  type_qn: string;
}

model VarRec {
  type_ref: TypeRef;
  conn_ph: string;
  expr_ck: int64;
}
