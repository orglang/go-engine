package commsem

type SemRefDS struct {
	CommID string `db:"comm_id"`
	CommRN int64  `db:"comm_rn"`
}
