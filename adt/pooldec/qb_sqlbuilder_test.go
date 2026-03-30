package pooldec

import (
	"fmt"
	"testing"
)

func TestSelectRecByQN(t *testing.T) {
	qb := newSQLBuilder()
	rec := decRecDS{}
	sql, _ := qb.insertRec(rec)
	fmt.Println(sql)
}
