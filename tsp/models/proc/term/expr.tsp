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

model CloseSpec { conn_ph: string; }
model WaitSpec { conn_ph: string; cont_expr: ExprSpec; }
model SendSpec { conn_ph: string; val_ph: string; }
model RecvSpec { conn_ph: string; bind_ph: string; cont_expr: ExprSpec; }
model LabSpec { conn_ph: string; pattern_qn: string; }
model CaseSpec { conn_ph: string; cont_bses: BranchSpec[]; }
model BranchSpec { pattern_qn: string; cont_expr: ExprSpec; }
model LinkSpec { conn_ph: string; cont_ph: string; }
model FwdSpec { conn_ph: string; cont_ph: string; }
model CallSpec {
  bind_ph: string;
  proc_qn: string;
  val_phs: string[];
  cont_expr: ExprSpec;
}
