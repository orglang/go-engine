package xactdef

const (
	descSems  string = "desc_sems sem"
	descBinds string = "desc_binds bind"
	xactDefs  string = "xact_defs def"
)

type queryBuilder interface {
	selectRecByQN() string
}
