package typedef

import "testing"

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(defRecDS{})
	t.Log(sql)
}

func TestSelectRecByQN(t *testing.T) {
	qb := newSQLBuilder()
	sql := qb.selectRecByQN()
	t.Log(sql)
}
