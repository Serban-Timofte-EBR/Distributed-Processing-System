package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Task struct {
	ID        string
	Operation string
	Args      []string
	Priority  string
}

var (
	highPriorityQueue   []Task
	mediumPriorityQueue []Task
	lowPriorityQueue    []Task
	taskResults         = make(map[string]net.Conn)
	workerPool          []net.Conn
	queueMutex          sync.Mutex
)

const (
	Reset   = "\033[0m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Red     = "\033[31m"
	Cyan    = "\033[36m"
	Magenta = "\033[35m"
)

func logInfo(color, label, msg string) {
	fmt.Printf("%s[SERVER][%s] %s%s\n", color, label, msg, Reset)
}

func addTaskToQueue(clientConn net.Conn, task string) {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	taskID := strconv.Itoa(len(taskResults) + 1)
	parts := strings.Fields(task)

	if len(parts) < 4 {
		logInfo(Red, "ERROR", "Invalid task format: "+task)
		return
	}

	operation := parts[0]
	args := parts[1 : len(parts)-1]
	priority := strings.ToUpper(parts[len(parts)-1])

	newTask := Task{
		ID:        taskID,
		Operation: operation,
		Args:      args,
		Priority:  priority,
	}

	switch priority {
	case "HIGH":
		highPriorityQueue = append(highPriorityQueue, newTask)
	case "MEDIUM":
		mediumPriorityQueue = append(mediumPriorityQueue, newTask)
	case "LOW":
		lowPriorityQueue = append(lowPriorityQueue, newTask)
	default:
		logInfo(Yellow, "WARNING", "Unknown priority, defaulting to MEDIUM")
		mediumPriorityQueue = append(mediumPriorityQueue, newTask)
	}

	taskResults[taskID] = clientConn

	logInfo(Cyan, "TASK_ADDED", fmt.Sprintf("%s %v [ID: %s] → [%s]", operation, args, taskID, priority))
}

func getNextTask() Task {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	if len(highPriorityQueue) > 0 {
		task := highPriorityQueue[0]
		highPriorityQueue = highPriorityQueue[1:]
		return task
	}
	if len(mediumPriorityQueue) > 0 {
		task := mediumPriorityQueue[0]
		mediumPriorityQueue = mediumPriorityQueue[1:]
		return task
	}
	if len(lowPriorityQueue) > 0 {
		task := lowPriorityQueue[0]
		lowPriorityQueue = lowPriorityQueue[1:]
		return task
	}
	return Task{}
}

func sendResultToClient(taskID, result string) {
	queueMutex.Lock()
	clientConn, exists := taskResults[taskID]
	queueMutex.Unlock()

	if exists {
		logInfo(Green, "RESULT", fmt.Sprintf("Sending result for Task ID %s → %s", taskID, result))
		_, err := clientConn.Write([]byte(result + "\n"))
		if err != nil {
			logInfo(Red, "ERROR", "Sending result failed: "+err.Error())
		}
	} else {
		logInfo(Red, "ERROR", "No client found for Task ID: "+taskID)
	}
}

func registerWorker(conn net.Conn) {
	queueMutex.Lock()
	workerPool = append(workerPool, conn)
	queueMutex.Unlock()
	logInfo(Magenta, "WORKER", "Registered: "+conn.RemoteAddr().String())
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			logInfo(Yellow, "CLOSED", "Connection closed: "+err.Error())
			conn.Close()
			return
		}
		message = strings.TrimSpace(message)

		switch {
		case message == "REGISTER_WORKER":
			registerWorker(conn)

		case message == "REQUEST_TASK":
			logInfo(Magenta, "WORKER", "Requested task")
			task := getNextTask()

			if task.ID != "" {
				logInfo(Cyan, "ASSIGN", fmt.Sprintf("Sending task [%s] to worker → %s %v", task.Priority, task.Operation, task.Args))
				taskStr := task.ID + " " + task.Operation + " " + strings.Join(task.Args, " ")
				_, err := conn.Write([]byte(taskStr + "\n"))
				if err != nil {
					logInfo(Red, "ERROR", "Sending task failed: "+err.Error())
					conn.Close()
					return
				}
			} else {
				_, _ = conn.Write([]byte("NO_TASK\n"))
			}

		case strings.HasPrefix(message, "RESULT"):
			parts := strings.Fields(message)
			if len(parts) < 3 {
				logInfo(Yellow, "WARNING", "Malformed result received")
				continue
			}
			taskID := parts[1]
			result := strings.Join(parts[2:], " ")
			logInfo(Green, "RESULT", "Received from worker: "+result)
			sendResultToClient(taskID, result)

		default:
			addTaskToQueue(conn, message)
		}
	}
}

func main() {
	logInfo(Cyan, "START", "Task Server started on port 50051")

	listener, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		logInfo(Red, "FATAL", "Error starting server: "+err.Error())
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logInfo(Red, "ERROR", "Accept connection failed: "+err.Error())
			continue
		}
		go handleConnection(conn)
	}
}
