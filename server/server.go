package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

var (
	taskQueue   []Task
	workerPool  []net.Conn
	taskResults = make(map[string]net.Conn)
	queueMutex  sync.Mutex
)

type Task struct {
	ID        string
	Operation string
	Args      []string
}

func addTaskToQueue(clientConn net.Conn, task string) {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	taskID := strconv.Itoa(len(taskResults) + 1) // TODO: Create a better ID system
	parts := strings.Fields(task)

	if len(parts) < 3 {
		fmt.Println("Error: Invalid task format.")
		return
	}

	newTask := Task{ID: taskID, Operation: parts[0], Args: parts[1:]}
	taskQueue = append(taskQueue, newTask)
	taskResults[taskID] = clientConn

	fmt.Println("Task added to queue:", newTask.Operation, newTask.Args, " (Task ID:", taskID, ")")
}

func sendResultToClient(taskID, result string) {
	queueMutex.Lock()
	clientConn, exists := taskResults[taskID]
	queueMutex.Unlock()

	if exists {
		fmt.Println("Sending result to client:", result)
		_, err := clientConn.Write([]byte(result + "\n"))
		if err != nil {
			fmt.Println("Error sending result to client:", err)
		}
	} else {
		fmt.Println("No client found for Task ID:", taskID)
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

			if task.ID != "" {
				fmt.Println("Assigning task to worker:", task.Operation, task.Args, "(Task ID:", task.ID, ")")
				taskStr := task.ID + " " + task.Operation + " " + strings.Join(task.Args, " ")
				_, err := conn.Write([]byte(taskStr + "\n"))
				if err != nil {
					fmt.Println("Error sending task to worker:", err)
					conn.Close()
					return
				}
			} else {
				fmt.Println("No tasks available.")
				_, _ = conn.Write([]byte("NO_TASK\n"))
			}
		} else if strings.HasPrefix(message, "RESULT") {
			parts := strings.Fields(message)
			if len(parts) < 3 {
				fmt.Println("Malformed result received.")
				continue
			}

			taskID := parts[1]
			result := strings.Join(parts[2:], " ")

			fmt.Println("Received result from worker:", result)
			sendResultToClient(taskID, result)
		} else {
			addTaskToQueue(conn, message)
		}
	}
}

func getNextTask() Task {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	if len(taskQueue) == 0 {
		return Task{}
	}

	task := taskQueue[0]
	taskQueue = taskQueue[1:]
	return task
}

func registerWorker(conn net.Conn) {
	queueMutex.Lock()
	workerPool = append(workerPool, conn)
	queueMutex.Unlock()
	fmt.Println("Worker registered:", conn.RemoteAddr().String())
}

func main() {
	fmt.Println("Task Server started on port 50051")

	listener, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
