package xactexp

import (
	"fmt"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(stateDS{})
	fmt.Println(sql)
}

func TestSelectRec(t *testing.T) {
	qb := newSQLBuilder()
	sql := qb.selectRecByVK()
	fmt.Println(sql)
}
