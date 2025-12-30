package type_test

import (
	"os"
	"testing"

	"orglang/orglang/adt/typedef"
)

var (
	api = typedef.NewAPI()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
