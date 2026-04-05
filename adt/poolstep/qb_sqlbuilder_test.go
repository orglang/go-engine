package poolstep

import (
	"fmt"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(StepRecDS{})
	fmt.Println(sql)
}
