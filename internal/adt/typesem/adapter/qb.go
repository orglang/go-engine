package adapter

type queryBuilder interface {
	updateRef(SemRefDS) (string, []any)
	selectRefByQN() string
}
