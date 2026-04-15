package commturn

import (
	"fmt"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(TurnRecDS{})
	fmt.Println(sql)
}
