package adapter

type SemRef struct {
	TermID string `db:"term_id"`
	DescID string `db:"desc_id"`
	DescRN int64  `db:"desc_rn"`
}
