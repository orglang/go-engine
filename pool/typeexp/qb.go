package typeexp

const (
	xactExps string = "pool_type_exps"
)

type queryBuilder interface {
	insertRec(stateDS) (string, []any)
	selectRecByVK() string
}
