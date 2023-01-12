package async

import (
	"fmt"
	"go.uber.org/zap"
)

type worker struct {
	id         int
	jobChannel chan Job
	quit       chan struct{}
	log        *zap.Logger
}

func NewWorker(id int, log *zap.Logger) *worker {
	return &worker{
		id:         id,
		jobChannel: make(chan Job),
		quit:       make(chan struct{}),
		log:        log,
	}
}

func (w *worker) Run(wq chan<- chan Job) {
	w.log.Info("worker init...", zap.Int("id", w.id))
	go func() {
		defer func() {
			panicErr := recover()
			if nil != panicErr {
				w.log.Error(fmt.Sprintf("【worker-%d】run painc, err: %s", w.id, panicErr))
				w.log.Info(fmt.Sprintf("【worker-%d】 recover worker...", w.id))
				w.Run(wq)
			}
		}()
		for {
			wq <- w.jobChannel
			select {
			case job := <-w.jobChannel:
				job.Do()
			case <-w.quit:
				w.log.Info(fmt.Sprintf("【worker-%d】worker quited", w.id))
				return
			}
		}
	}()
}

func (w *worker) Close() {
	go func() {
		w.quit <- struct{}{}
	}()
}
