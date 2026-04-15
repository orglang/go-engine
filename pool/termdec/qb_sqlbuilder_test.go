package termdec

import (
	"fmt"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(decRecDS{})
	fmt.Println(sql)
}

func TestSelectRecByQN(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.selectRecByQN("qn")
	fmt.Println(sql)
}
