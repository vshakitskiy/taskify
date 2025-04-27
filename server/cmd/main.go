package main

import (
	"app.server/internal/app"
	"app.server/internal/config"
	"github.com/fatih/color"
)

func main() {
	if err := config.InitConfig(); err != nil {
		color.Red("Unable to init config: %s", err)
		return
	}

	a := app.NewApp()

	if err := a.Run(); err != nil {
		panic(err)
	}
}

// func main() {
// 	mq, err := mq.NewMQ("amqp://admeanie:shabi@localhost:5672/", 3)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer mq.Methods.Close()

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	taskResults := make(map[string]chan string, 50)
// 	go func() {
// 		for del := range mq.ResultsCh {
// 			fmt.Println(del.CorrelationId, string(del.Body))
// 			if _, ok := taskResults[del.CorrelationId]; !ok {
// 				fmt.Println("not found")
// 				del.Ack(false)
// 				continue
// 			}
// 			taskResults[del.CorrelationId] <- string(del.Body)
// 			del.Ack(false)
// 		}
// 	}()

// 	go func() {
// 		for {
// 			fmt.Println("Current task results len", len(taskResults))

// 			time.Sleep(1 * time.Second)

// 			task := models.NewTask("blabla", 3000)
// 			id := uuid.New().String()
// 			replyTo := "results"
// 			taskResults[id] = make(chan string, 1)

// 			mq.Methods.Publish(
// 				ctx,
// 				mq.TaskQueue,
// 				task.ToJSON(),
// 				&id,
// 				&replyTo,
// 			)

// 			res := <-taskResults[id]
// 			close(taskResults[id])
// 			fmt.Println("got from channel", res)
// 			delete(taskResults, id)
// 		}
// 	}()

// 	shutdownChan := make(chan os.Signal, 1)
// 	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

// 	<-shutdownChan
// }
