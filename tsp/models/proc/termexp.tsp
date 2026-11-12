namespace OrgLang.Engine.Proc;

/** Process term expression specification (go-sdk/proc/termexp). */
model TermExprSpec {
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

model CloseSpec { comm_ph: string; }
model WaitSpec { comm_ph: string; cont_es: TermExprSpec; }
model SendSpec { comm_ph: string; val_ph: string; }
model RecvSpec { comm_ph: string; bind_ph: string; cont_es: TermExprSpec; }
model LabSpec { comm_ph: string; pattern_qn: string; }
model CaseSpec { comm_ph: string; cont_bses: BranchSpec[]; }
model BranchSpec { pattern_qn: string; cont_es: TermExprSpec; }
model CallSpec {
  bind_ph: string;
  proc_qn: string;
  val_phs: string[];
  cont_es: TermExprSpec;
}
model FwdSpec { comm_ph: string; cont_ph: string; }
