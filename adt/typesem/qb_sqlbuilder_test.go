package typesem

import (
	"fmt"
	"testing"
)

func TestUpdateRef(t *testing.T) {
	qb := newSQLBuilder("types", "descs")
	sql, _ := qb.updateRef(SemRefDS{})
	fmt.Println(sql)
}
