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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err)
		conn.Close()
		return
	}
	message = strings.TrimSpace(message)

	if message == "REGISTER_WORKER" {
		registerWorker(conn)
	} else {
		result := assignTaskToWorker(message)
		conn.Write([]byte(result + "\n"))
	}
}

func registerWorker(conn net.Conn) {
	workerMutex.Lock()
	workers = append(workers, conn)
	workerMutex.Unlock()
	fmt.Println("Worker registered:", conn.RemoteAddr())

	// Keep the worker connection open to continuously listen for tasks
	go workerListener(conn)
}

func assignTaskToWorker(task string) string {
	workerMutex.Lock()
	defer workerMutex.Unlock()

	if len(workers) == 0 {
		return "No workers available"
	}

	worker := workers[0] // Get the first worker (round-robin)
	workers = append(workers[1:], worker)

	_, err := worker.Write([]byte(task + "\n"))
	if err != nil {
		fmt.Println("Error sending task to worker:", err)
		worker.Close()
		return "Worker communication failed"
	}

	response, err := bufio.NewReader(worker).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response from worker:", err)
		worker.Close()
		return "Worker response failed"
	}

	return strings.TrimSpace(response)
}

func workerListener(worker net.Conn) {
	defer worker.Close()

	reader := bufio.NewReader(worker)
	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Worker disconnected:", worker.RemoteAddr())
			removeWorker(worker)
			return
		}
	}
}

func removeWorker(worker net.Conn) {
	workerMutex.Lock()
	defer workerMutex.Unlock()

	for i, w := range workers {
		if w == worker {
			workers = append(workers[:i], workers[i+1:]...)
			break
		}
	}
}
