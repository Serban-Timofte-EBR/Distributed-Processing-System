package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
)

func logLine(color string, msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s[%s] %s%s\n", color, timestamp, msg, colorReset)
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:50051")
	if err != nil {
		logLine(colorRed, "Eroare conectare la server: "+err.Error())
		return
	}
	defer conn.Close()

	logLine(colorGreen, "Task Worker connected.")
	_, err = conn.Write([]byte("REGISTER_WORKER\n"))
	if err != nil {
		logLine(colorRed, "Eroare trimitere mesaj REGISTER_WORKER: "+err.Error())
		return
	}

	reader := bufio.NewReader(conn)

	for {
		logLine(colorCyan, "Worker requesting task from server...")
		_, err = conn.Write([]byte("REQUEST_TASK\n"))
		if err != nil {
			logLine(colorRed, "Eroare trimitere mesaj REQUEST_TASK: "+err.Error())
			return
		}

		taskLine, err := reader.ReadString('\n')
		if err != nil {
			logLine(colorRed, "Worker connection lost: "+err.Error())
			return
		}

		taskParts := strings.Fields(strings.TrimSpace(taskLine))
		if len(taskParts) < 3 {
			logLine(colorYellow, "Received malformed task. Ignoring.")
			time.Sleep(time.Second)
			continue
		}

		taskID := taskParts[0]
		operation := taskParts[1]
		args := taskParts[2:]

		logLine(colorCyan, fmt.Sprintf("Executing Task: %s %s", operation, strings.Join(args, " ")))
		result := executeTask(operation, args)

		response := "RESULT " + taskID + " " + result
		logLine(colorGreen, "Sending result back to server: "+result)

		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			logLine(colorRed, "Eroare trimitere rezultat: "+err.Error())
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
