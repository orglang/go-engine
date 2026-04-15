package semterm

type TermRefDS struct {
	TermID string `db:"impl_id"`
	TermRN int64  `db:"impl_rn"`
}
