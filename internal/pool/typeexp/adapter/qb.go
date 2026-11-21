package adapter

type QueryBuilder interface {
	InsertRec(State) (string, []any)
	SelectRecByCK() string
}

const (
	xactExps string = "pool_type_exps"
)
