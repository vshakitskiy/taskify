package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"shared/models"
	"sync"
	"syscall"
	"time"
)

func worker(id int, taskQueue <-chan *models.Task, resultChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d starting\n", id)

	for task := range taskQueue {
		fmt.Printf("Worker %d: Received %s (ID: %s)\n", id, task.Message, task.Id)
		time.Sleep(task.Timeout)
		fmt.Printf("Worker %d: Finished %s (ID: %s)\n", id, task.Message, task.Id)

		resultChan <- fmt.Sprintf("%s: %s", task.Id, task.Message)
	}

	fmt.Printf("Worker %d stopping\n", id)
}

func taskGenerator(taskQueue chan<- *models.Task, stopChan <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	counter := 1

	fmt.Println("Task generator starting")

	for {
		select {
		case <-ticker.C:
			task := models.NewTask(
				fmt.Sprintf("task %d", counter),
				time.Duration(rand.Intn(10))*time.Second,
			)

			fmt.Printf("Generator: Sending %s (ID: %s, Timeout: %v)\n", task.Message, task.Id, task.Timeout)

			select {
			case taskQueue <- task:
				counter++
			default:
				fmt.Println("Task queue full")
			}
		case <-stopChan:
			fmt.Println("Task generator stopping...")
			return
		}
	}
}

func resultProcessor(resultChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Result processor starting")
	for result := range resultChan {
		fmt.Println(result)
	}
	fmt.Println("Result processor stopping")
}

func main() {
	// TODO: make better workers implementation
	// TODO: connect with RabbitMQ

	taskQueue := make(chan *models.Task, 10)
	resultChan := make(chan string)
	stopGeneratorChan := make(chan struct{})

	var generatorWg sync.WaitGroup
	var workerWg sync.WaitGroup
	var resultWg sync.WaitGroup

	generatorWg.Add(1)
	go taskGenerator(taskQueue, stopGeneratorChan, &generatorWg)

	numWorkers := 3
	workerWg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go worker(i, taskQueue, resultChan, &workerWg)
	}

	resultWg.Add(1)
	go resultProcessor(resultChan, &resultWg)

	fmt.Println("Application started. Press Ctrl+C to exit.")
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	<-shutdownChan
	fmt.Println("\nShutdown signal received...")

	fmt.Println("Signaling generator to stop...")
	close(stopGeneratorChan)

	fmt.Println("Waiting for generator to finish...")
	generatorWg.Wait()
	fmt.Println("Generator finished.")

	fmt.Println("Closing task queue...")
	close(taskQueue)

	fmt.Println("Waiting for workers to finish...")
	workerWg.Wait()
	fmt.Println("Workers finished.")

	fmt.Println("Closing result channel...")
	close(resultChan)

	fmt.Println("Waiting for result processor to finish...")
	resultWg.Wait()
	fmt.Println("Result processor finished.")

	fmt.Println("Application shutdown complete.")
}
