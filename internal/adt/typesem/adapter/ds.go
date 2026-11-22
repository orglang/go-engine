package adapter

type SemRef struct {
	TypeID string `db:"type_id"`
	DescID string `db:"desc_id"`
	DescRN int64  `db:"desc_rn"`
}
