package poolconn

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	qb := newSQLBuilder()
	rec := connRecDS{CommID: "id-1", CommRN: 1}
	sql, args := qb.insertRec(rec)
	fmt.Println(sql)
	fmt.Println(args)
}
