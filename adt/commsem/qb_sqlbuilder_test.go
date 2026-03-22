package commsem

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	qb := newSQLBuilder()
	rec := semRecDS{CommID: "id-1", CommRN: 1, Kind: 2}
	sql, args := qb.insertRec(rec)
	fmt.Println(sql)
	fmt.Println(args)
}
