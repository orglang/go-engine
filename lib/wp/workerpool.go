package wp

import (
	"github.com/gammazero/workerpool"
)

func newWorkerPool() *workerpool.WorkerPool {
	return workerpool.New(10)
}
