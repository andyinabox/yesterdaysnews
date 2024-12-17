package batchwg

import "sync"

type WaitGroup struct {
	wg        sync.WaitGroup
	batchSize int
	current   int
}

func New(batchSize int) *WaitGroup {
	return &WaitGroup{
		batchSize: batchSize,
	}
}

func (w *WaitGroup) Add(n int) {
	w.wg.Add(n)
	w.current += n
}

func (w *WaitGroup) Wait() {
	if w.current >= w.batchSize {
		w.wg.Wait()
		w.current = 0
	}
}

func (w *WaitGroup) Done() {
	w.wg.Done()
}
