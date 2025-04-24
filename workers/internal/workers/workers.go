package workers

import (
	"context"
	"fmt"
	"shared/models"
	"sync"
	"time"

	"app.workers/internal/mq"
	"github.com/fatih/color"
)

type Workers struct {
	mq     *mq.RabbitMQ
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	amount int
}

func NewWorkers(mq *mq.RabbitMQ, amount int) *Workers {
	ctx, cancel := context.WithCancel(context.Background())
	return &Workers{
		mq:     mq,
		wg:     &sync.WaitGroup{},
		ctx:    ctx,
		cancel: cancel,
		amount: amount,
	}
}

func (w *Workers) Run() {
	for i := 0; i < w.amount; i++ {
		w.wg.Add(1)
		go worker(i, w.wg, w.ctx, w.mq)
	}
}

func (w *Workers) Stop() {
	w.cancel()
	w.wg.Wait()
}

func worker(
	id int,
	wg *sync.WaitGroup,
	ctx context.Context,
	mq *mq.RabbitMQ,
) {
	color.Blue("Starting worker %d...", id)

	for {
		select {
		case del := <-mq.TasksCh:
			task, err := models.TaskFromJSON(del.Body)
			if err != nil {
				color.Red("Failed to unmarshal task: %s", err)

				res := models.NewErrorTaskResponse("ErrInvalidPayload")
				err = mq.Methods.Publish(
					ctx,
					mq.ResultsQueue,
					res.ToJSON(),
					&del.CorrelationId,
					&del.ReplyTo,
				)
				if err != nil {
					color.Red("Failed to publish a message: %s", err)
					del.Nack(false, false)
					continue
				}

				del.Ack(false)
				continue
			}

			rand := time.Duration(task.Timeout) * time.Millisecond
			time.Sleep(rand)

			fmt.Printf("Worker %d: Received %s (ID: %s)\n", id, task.Message, del.CorrelationId)

			res := models.NewSuccessTaskResponse(task.Message)

			err = mq.Methods.Publish(
				ctx,
				mq.ResultsQueue,
				res.ToJSON(),
				&del.CorrelationId,
				&del.ReplyTo,
			)
			if err != nil {
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
