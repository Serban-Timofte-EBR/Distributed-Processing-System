package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const (
	serverAddress = "localhost:50051"
)

func main() {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Task Worker connected.")

	// Register with the server
	_, err = conn.Write([]byte("REGISTER_WORKER\n"))
	if err != nil {
		fmt.Println("Error registering worker:", err)
		return
	}

	// Continuously listen for tasks
	reader := bufio.NewReader(conn)
	for {
		task, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Worker connection lost:", err)
			return
		}

		task = strings.TrimSpace(task)
		fmt.Println("Executing Task:", task)

		result := executeTask(task)
		_, err = conn.Write([]byte(result + "\n"))
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}
}

func executeTask(task string) string {
	parts := strings.Split(task, " ")
	if len(parts) < 3 {
		return "Invalid task format"
	}

	if parts[0] == "Multiply" {
		num1 := parseInt(parts[1])
		num2 := parseInt(parts[2])
		return fmt.Sprintf("Result: %d", num1*num2)
	}
	return "Unknown Task"
}

func parseInt(s string) int {
	var num int
	fmt.Sscanf(s, "%d", &num)
	return num
}
