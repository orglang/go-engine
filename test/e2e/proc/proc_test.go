package sig_test

import (
	"os"
	"testing"

	"orglang/orglang/adt/procdecl"
)

var (
	api = procdecl.NewAPI()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
