import "../type.tsp";

namespace OrgLang.Any.Comm;

model VarSpec {
  comm_ph: string;
  type_qn: string;
}

model VarRec {
  type_ref: TypeRef;
  comm_ph: string;
  expr_ck: int64;
}
