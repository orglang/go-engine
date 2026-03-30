package pooldec

const (
	poolDecs string = "pool_decs dec"
)

type queryBuilder interface {
	insertRec(decRecDS) (string, []any)
}
