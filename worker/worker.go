package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
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

	_, err = conn.Write([]byte("REGISTER_WORKER\n"))
	if err != nil {
		fmt.Println("Error registering worker:", err)
		return
	}

	reader := bufio.NewReader(conn)

	for {
		fmt.Println("Worker requesting task from server...")

		_, err = conn.Write([]byte("REQUEST_TASK\n"))
		if err != nil {
			fmt.Println("Error requesting task:", err)
			return
		}

		task, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Worker connection lost:", err)
			return
		}

		task = strings.TrimSpace(task)

		fmt.Println("Worker received from server:", task) // ✅ Debugging

		if !isValidTask(task) {
			fmt.Println("Received malformed task. Ignoring.")
			time.Sleep(time.Second)
			continue
		}

		fmt.Println("Executing Task:", task)
		result := executeTask(task)

		fmt.Println("Sending result back to server:", result)

		_, err = conn.Write([]byte(result + "\n"))
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}
}

func isValidTask(task string) bool {
	parts := strings.Fields(task)
	if len(parts) < 3 {
		return false
	}
	return true
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
