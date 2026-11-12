namespace OrgLang.Engine.Common;

model TermRef {
  term_id: string;
  term_rn: int64;
}

model TypeRef {
  type_id: string;
  type_rn: int64;
}

/** Reference to a computation (the Go wire prefix is comp_). */
model ComputationRef {
  comp_id: string;
  comp_rn: int64;
}

model TypeLinkSpec { type_qn: string; }

model TermVarSpec {
  chnl_ph: string;
  type_qn: string;
}

model CompVarSpec {
  chnl_ph: string;
  term_qn: string;
}

model TermVarRec {
  type_ref: TypeRef;
  chnl_ph: string;
  exp_vk: int64;
}
