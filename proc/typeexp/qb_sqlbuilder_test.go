package typeexp

import "testing"

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(stateDS{})
	t.Log(sql)
}

func TestSelectRec(t *testing.T) {
	qb := newSQLBuilder()
	sql := qb.selectRecByVK()
	fmt.Println(sql)
}
