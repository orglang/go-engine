namespace OrgLang.Proc.Term;

/** Process term expression specification. Exactly one variant is serialized inline with its kind discriminator. */
@discriminated(#{ envelope: "none" })
union ExprSpec {
  close: CloseSpec,
  wait: WaitSpec,
  send: SendSpec,
  recv: RecvSpec,
  lab: LabSpec,
  case: CaseSpec,
  link: LinkSpec,
  fwd: FwdSpec,
  call: CallSpec,
}

model CloseSpec {
  comm_ph: string;
}
model WaitSpec {
  comm_ph: string;
  cont_expr: ExprSpec;
}
model SendSpec {
  comm_ph: string;
  val_ph: string;
}
model RecvSpec {
  comm_ph: string;
  bind_ph: string;
  cont_expr: ExprSpec;
}
model LabSpec {
  comm_ph: string;
  pattern_qn: string;
}
model CaseSpec {
  comm_ph: string;
  cont_bses: BranchSpec[];
}
model BranchSpec {
  pattern_qn: string;
  cont_expr: ExprSpec;
}
model LinkSpec {
  proc_qn: string;
  x: string;
  ys: string[];
}
model FwdSpec {
  comm_ph: string;
  cont_ph: string;
}
model CallSpec {
  bind_ph: string;
  proc_qn: string;
  val_phs: string[];
  cont_expr: ExprSpec;
}
