package sig_test

import (
	"os"
	"testing"

	"orglang/orglang/adt/procdec"
)

var (
	api = procdec.NewAPI()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
