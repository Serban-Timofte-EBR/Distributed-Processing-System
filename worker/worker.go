package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)


var (
	host    = flag.String("host", "127.0.0.1", "Server host")
	port    = flag.Int("port", 50051, "Server port")
	logFile = flag.String("log", "app-logs/worker.log", "Path to worker log file")
)

func init() {
	f, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime)
}

func logLine(msg string) {
	log.Printf("%s", msg)
}

func executeTask(op string, args []string) string {
	x, y := toInt(args[0]), toInt(args[1])
	switch op {
	case "Add":
		return fmt.Sprintf("%d + %d = %d", x, y, x+y)
	case "Minus":
		return fmt.Sprintf("%d - %d = %d", x, y, x-y)
	case "Multiply":
		return fmt.Sprintf("%d * %d = %d", x, y, x*y)
	case "Subtract":
		return fmt.Sprintf("%d / %d = %d", x, y, x/y)
	case "Power":
		return fmt.Sprintf("%d ^ %d = %d", x, y, pow(x, y))
	case "Mod":
		return fmt.Sprintf("%d %% %d = %d", x, y, x%y)
	}
	return "INVALID_TASK"
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func pow(a, b int) int {
	result := 1
	for i := 0; i < b; i++ {
		result *= a
	}
	return result
}

func main() {
	flag.Parse()
	address := fmt.Sprintf("%s:%d", *host, *port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		logLine("Error connecting to server: " + err.Error())
		return
	}
	defer conn.Close()

	logLine("Worker connected to " + address)
	conn.Write([]byte("REGISTER_WORKER\n"))

	r := bufio.NewReader(conn)
	for {
		logLine("Requesting task...")
		conn.Write([]byte("REQUEST_TASK\n"))
		line, err := r.ReadString('\n')
		if err != nil {
			logLine("Connection lost: " + err.Error())
			return
		}

		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) < 3 {
			time.Sleep(time.Second)
			continue
		}
		id, op, args := parts[0], parts[1], parts[2:]
		logLine(fmt.Sprintf("Executing [%s]: %v", id, parts[1:]))
		res := executeTask(op, args)
		logLine("Sending result: " + res)
		conn.Write([]byte(fmt.Sprintf("RESULT %s %s\n", id, res)))
	}
}