namespace OrgLang.Engine;

model TermRef {
  term_id: string;
  term_rn: int64;
}

model TypeRef {
  type_id: string;
  type_rn: int64;
}

/** Reference to a computation, not a generic component. */
model ComputationRef {
  comp_id: string;
  comp_rn: int64;
}

model TermVarSpec {
  chnl_ph: string;
  type_qn: string;
}

model ComputationVarSpec {
  chnl_ph: string;
  term_qn: string;
}

model TermVarRec {
  type_ref: TypeRef;
  chnl_ph: string;
  exp_vk: int64;
}
