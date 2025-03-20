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

		go handleClient(conn)
	}
}

func handleClient(clientConn net.Conn) {
	defer clientConn.Close()

	reader := bufio.NewReader(clientConn)
	taskRequest, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from client:", err)
		return
	}

	taskRequest = strings.TrimSpace(taskRequest)
	fmt.Println("Received Task:", taskRequest)

	// Assign task to an available worker
	result := assignTaskToWorker(taskRequest)
	clientConn.Write([]byte(result + "\n"))
}

func assignTaskToWorker(task string) string {
	workerMutex.Lock()
	defer workerMutex.Unlock()

	if len(workers) == 0 {
		return "No workers available"
	}

	worker := workers[0]
	workers = append(workers[1:], worker) // Round-robin worker selection

	_, err := worker.Write([]byte(task + "\n"))
	if err != nil {
		fmt.Println("Error sending task to worker:", err)
		return "Worker communication failed"
	}

	response, _ := bufio.NewReader(worker).ReadString('\n')
	return strings.TrimSpace(response)
}

func registerWorker(conn net.Conn) {
	workerMutex.Lock()
	workers = append(workers, conn)
	workerMutex.Unlock()
}
