package connsem

type SemRefDS struct {
	ConnID string `db:"comm_id"`
	ConnRN int64  `db:"comm_rn"`
}
