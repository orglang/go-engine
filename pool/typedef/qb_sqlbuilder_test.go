package typedef

import (
	"fmt"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(defRecDS{})
	fmt.Println(sql)
}

func TestSelectRecByQN(t *testing.T) {
	qb := newSQLBuilder()
	sql := qb.selectRecByQN()
	fmt.Println(sql)
}
