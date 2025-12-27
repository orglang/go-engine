package sig_test

import (
	"os"
	"testing"

	procdec "smecalculus/rolevod/app/proc/dec"
)

var (
	api = procdec.NewAPI()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
