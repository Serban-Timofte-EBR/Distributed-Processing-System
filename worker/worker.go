package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:50051")
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

		taskLine, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Worker connection lost:", err)
			return
		}

		taskParts := strings.Fields(strings.TrimSpace(taskLine))
		if len(taskParts) < 3 {
			fmt.Println("Received malformed task. Ignoring.")
			time.Sleep(time.Second)
			continue
		}

		taskID := taskParts[0]
		operation := taskParts[1]
		args := taskParts[2:]

		fmt.Println("Executing Task:", operation, args)
		result := executeTask(operation, args)

		response := "RESULT " + taskID + " " + result
		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}
}

func executeTask(operation string, args []string) string {
	if operation == "Multiply" && len(args) == 2 {
		return args[0] + " * " + args[1] + " = " + fmt.Sprintf("%d", multiply(args[0], args[1]))
	}
	return "INVALID_TASK"
}

func multiply(a, b string) int {
	var x, y int
	fmt.Sscanf(a, "%d", &x)
	fmt.Sscanf(b, "%d", &y)
	return x * y
}
