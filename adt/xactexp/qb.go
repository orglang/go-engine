package xactexp

const (
	xactExps string = "xact_exps"
)

type queryBuilder interface {
	insertRec(stateDS) (string, []any)
	selectRecByVK() string
}
