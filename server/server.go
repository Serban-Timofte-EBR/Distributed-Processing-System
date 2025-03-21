package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

const (
	serverPort = "50051"
)

var workers []net.Conn
var workerMutex sync.Mutex
var taskQueue []string
var queueMutex sync.Mutex

func main() {
	listener, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Task Server started on port", serverPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed:", err)
			conn.Close()
			return
		}
		message = strings.TrimSpace(message)

		if message == "REGISTER_WORKER" {
			registerWorker(conn)
		} else if message == "REQUEST_TASK" {
			fmt.Println("Worker requested a task.")

			task := getNextTask()

			_, err := conn.Write([]byte(task + "\n"))
			if err != nil {
				fmt.Println("Error sending task to worker:", err)
				conn.Close()
				return
			}

			fmt.Println("Sent task to worker:", task)
		} else {
			addTaskToQueue(message)
		}
	}
}

func registerWorker(conn net.Conn) {
	workerMutex.Lock()
	workers = append(workers, conn)
	workerMutex.Unlock()
	fmt.Println("Worker registered:", conn.RemoteAddr())
}

func addTaskToQueue(task string) {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	if strings.TrimSpace(task) == "" {
		fmt.Println("Error: Received empty task. Ignoring.")
		return
	}

	fmt.Println("Task added to queue:", task)
	taskQueue = append(taskQueue, task)
}

func getNextTask() string {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	fmt.Println("Worker requested a task. Current queue size:", len(taskQueue))

	if len(taskQueue) == 0 {
		fmt.Println("Worker requested a task, but queue is empty.")
		return "NO_TASK"
	}

	task := taskQueue[0]
	taskQueue = taskQueue[1:]

	fmt.Println("Assigning task to worker:", task)
	if task == "" {
		fmt.Println("Error: Task is empty! Assigning NO_TASK.")
		return "NO_TASK"
	}

	return task
}
