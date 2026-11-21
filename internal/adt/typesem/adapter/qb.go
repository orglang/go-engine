package adapter

type queryBuilder interface {
	updateRef(SemRef) (string, []any)
	selectRefByQN() string
}
