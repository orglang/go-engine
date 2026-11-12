namespace OrgLang.Engine.Comp;

/** Computation models; the Go wire prefix "comp" means computation. */
model Ref {
  comp_id: string;
  comp_rn: int64;
}

model VarSpec {
  chnl_ph: string;
  term_qn: string;
}

model Snap {
  comp_ref: Ref;
}
