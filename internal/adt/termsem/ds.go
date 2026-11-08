package termsem

type SemRefDS struct {
	TermID string `db:"term_id"`
	TermRN int64  `db:"term_rn"`
}
