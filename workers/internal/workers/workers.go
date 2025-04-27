package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"app.shared/models"

	"app.workers/internal/mq"
	"github.com/fatih/color"
)

type Workers struct {
	mq            *mq.RabbitMQ
	ctx           context.Context
	cancel        context.CancelFunc
	publishCtx    context.Context
	publishCancel context.CancelFunc
	wg            *sync.WaitGroup
	amount        int
}

func NewWorkers(mq *mq.RabbitMQ, amount int) *Workers {
	ctx, cancel := context.WithCancel(context.Background())
	publishCtx, publishCancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	return &Workers{
		mq:            mq,
		wg:            &sync.WaitGroup{},
		ctx:           ctx,
		cancel:        cancel,
		publishCtx:    publishCtx,
		publishCancel: publishCancel,
		amount:        amount,
	}
}

func (w *Workers) Run() {
	for i := 0; i < w.amount; i++ {
		w.wg.Add(1)
		go worker(i, w.wg, w.ctx, w.publishCtx, w.mq)
	}
}

func (w *Workers) Stop() {
	w.cancel()
	w.wg.Wait()
	w.publishCancel()
}

func worker(
	id int,
	wg *sync.WaitGroup,
	ctx context.Context,
	pubCtx context.Context,
	mq *mq.RabbitMQ,
) {
	color.Blue("Starting worker %d...", id)

	for {
		select {
		case del := <-mq.TasksCh:
			task, err := models.TaskFromJSON(del.Body)
			if err != nil {
				color.Red("[%d] Failed to unmarshal task %s: %s", id, del.CorrelationId, err)

				res := models.NewErrorTaskResponse(id, "ErrInvalidPayload")
				if err = mq.Methods.Publish(
					pubCtx,
					mq.ResultsQueue,
					res.ToJSON(),
					&del.CorrelationId,
					&del.ReplyTo,
				); err != nil {
					color.Red(
						"[%d] Failed to publish a message (ID: %s): %s",
						id,
						del.CorrelationId,
						err,
					)
					del.Nack(false, false)
					continue
				}

				del.Ack(false)
				continue
			}
			fmt.Printf(
				"[%d]: %s - %dms (ID: %s)\n",
				id,
				task.Message,
				task.Timeout,
				del.CorrelationId,
			)

			rand := time.Duration(task.Timeout) * time.Second
			time.Sleep(rand)
			res := models.NewSuccessTaskResponse(id, task.Message)

			if err = mq.Methods.Publish(
				pubCtx,
				mq.ResultsQueue,
				res.ToJSON(),
				&del.CorrelationId,
				&del.ReplyTo,
			); err != nil {
				color.Red("Failed to publish a message: %s", err)
				del.Nack(false, false)
				continue
			}

			del.Ack(false)
		case <-ctx.Done():
			color.Blue("Stopping worker %d...", id)
			wg.Done()
			return
		}
	}
}
