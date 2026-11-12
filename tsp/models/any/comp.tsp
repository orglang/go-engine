namespace OrgLang.Engine.Any;

/** Shared computation identity used by process and pool computations. */
model CompRef {
  comp_id: string;
  comp_rn: int64;
}
