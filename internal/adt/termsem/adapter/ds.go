package adapter

type SemRef struct {
	TermID string `db:"term_id"`
	TermRN int64  `db:"term_rn"`
}
