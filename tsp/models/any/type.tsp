namespace OrgLang.Any;

/** Shared type identity. */
model TypeRef {
  type_id: string;
  type_rn: int64;
}

namespace OrgLang.Any.Type;

model LinkSpec {
  type_qn: string;
}
