package compsem

import (
	"fmt"
	"testing"
)

func TestUpdateRef(t *testing.T) {
	qb := newSQLBuilder("foo")
	sql, _ := qb.updateRef(SemRefDS{})
	fmt.Println(sql)
}
