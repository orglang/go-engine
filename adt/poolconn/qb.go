package poolconn

const (
	poolConns string = "pool_conns conn"
)

type queryBuilder interface {
	insertRec(connRecDS) (string, []any)
}
