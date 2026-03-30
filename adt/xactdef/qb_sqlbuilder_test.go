package xactdef

import (
	"fmt"
	"testing"
)

func TestSelectRecByQN(t *testing.T) {
	qb := newSQLBuilder()
	sql := qb.selectRecByQN()
	fmt.Println(sql)
}
