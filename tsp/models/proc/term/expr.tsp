namespace OrgLang.Proc.Term;

/** Process term expression specification. */
model ExprSpec {
  kind: "close" | "wait" | "send" | "recv" | "lab" | "case" | "link" | "fwd" | "call";
  close?: CloseSpec;
  wait?: WaitSpec;
  send?: SendSpec;
  recv?: RecvSpec;
  lab?: LabSpec;
  case?: CaseSpec;
  call?: CallSpec;
  fwd?: FwdSpec;
}

model CloseSpec { conn_ph: string; }
model WaitSpec { conn_ph: string; cont_exprr: ExprSpec; }
model SendSpec { conn_ph: string; val_ph: string; }
model RecvSpec { conn_ph: string; bind_ph: string; cont_exprr: ExprSpec; }
model LabSpec { conn_ph: string; pattern_qn: string; }
model CaseSpec { conn_ph: string; cont_bses: BranchSpec[]; }
model BranchSpec { pattern_qn: string; cont_exprr: ExprSpec; }
model CallSpec {
  bind_ph: string;
  proc_qn: string;
  val_phs: string[];
  cont_exprr: ExprSpec;
}
model FwdSpec { conn_ph: string; cont_ph: string; }
