package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const (
	serverAddress = "localhost:50051"
	workerPort    = "50052"
)

func main() {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Task Worker connected.")

	for {
		// Read task from server
		task, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading task:", err)
			return
		}

		task = strings.TrimSpace(task)
		fmt.Println("Executing Task:", task)

		result := executeTask(task)
		conn.Write([]byte(result + "\n"))
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
