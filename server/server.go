package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	taskQueue   []Task
	workerPool  []net.Conn
	taskResults = make(map[string]net.Conn)
	queueMutex  sync.Mutex
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[34m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorRed    = "\033[31m"
)

func logLine(color string, msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s[%s] %s%s\n", color, timestamp, msg, colorReset)
}

type Task struct {
	ID        string
	Operation string
	Args      []string
}

func addTaskToQueue(clientConn net.Conn, task string) {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	taskID := strconv.Itoa(len(taskResults) + 1)
	parts := strings.Fields(task)

	if len(parts) < 3 {
		logLine(colorRed, "Error: Invalid task format.")
		return
	}

	newTask := Task{ID: taskID, Operation: parts[0], Args: parts[1:]}
	taskQueue = append(taskQueue, newTask)
	taskResults[taskID] = clientConn

	logLine(colorBlue, fmt.Sprintf("Task added to queue: %s %v (Task ID: %s)", newTask.Operation, newTask.Args, taskID))
}

func sendResultToClient(taskID, result string) {
	queueMutex.Lock()
	clientConn, exists := taskResults[taskID]
	queueMutex.Unlock()

	if exists {
		logLine(colorCyan, fmt.Sprintf("Sending result to client: %s", result))
		_, err := clientConn.Write([]byte(result + "\n"))
		if err != nil {
			logLine(colorRed, "Error sending result to client: "+err.Error())
		}
	} else {
		logLine(colorRed, "No client found for Task ID: "+taskID)
	}
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			logLine(colorRed, "Connection closed: "+err.Error())
			conn.Close()
			return
		}
		message = strings.TrimSpace(message)

		switch {
		case message == "REGISTER_WORKER":
			registerWorker(conn)

		case message == "REQUEST_TASK":
			logLine(colorGreen, "Worker requested a task.")
			task := getNextTask()

			if task.ID != "" {
				logLine(colorYellow, fmt.Sprintf("Assigning task to worker: %s %v (Task ID: %s)", task.Operation, task.Args, task.ID))
				taskStr := task.ID + " " + task.Operation + " " + strings.Join(task.Args, " ")
				_, err := conn.Write([]byte(taskStr + "\n"))
				if err != nil {
					logLine(colorRed, "Error sending task to worker: "+err.Error())
					conn.Close()
					return
				}
			} else {
				logLine(colorYellow, "No tasks available. Sending NO_TASK.")
				_, _ = conn.Write([]byte("NO_TASK\n"))
			}

		case strings.HasPrefix(message, "RESULT"):
			parts := strings.Fields(message)
			if len(parts) < 3 {
				logLine(colorRed, "Malformed result received.")
				continue
			}

			taskID := parts[1]
			result := strings.Join(parts[2:], " ")
			logLine(colorCyan, fmt.Sprintf("Received result from worker (Task ID: %s): %s", taskID, result))
			sendResultToClient(taskID, result)

		default:
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
	logLine(colorGreen, "Worker registered: "+conn.RemoteAddr().String())
}

func main() {
	logLine(colorGreen, "Task Server started on port 50051")

	listener, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		logLine(colorRed, "Error starting server: "+err.Error())
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logLine(colorRed, "Error accepting connection: "+err.Error())
			continue
		}
		go handleConnection(conn)
	}
}
