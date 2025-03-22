package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
	colorCyan  = "\033[36m"
)

func logLine(color string, label string, msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s[%s] [%s] %s%s\n", color, timestamp, label, msg, colorReset)
}

func generateTasks() []string {
	var tasks []string
	for i := 1; i <= 50; i++ {
		tasks = append(tasks, fmt.Sprintf("Multiply %d %d", i, i+1))
	}
	return tasks
}

func sendTask(task string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "127.0.0.1:50051")
	if err != nil {
		logLine(colorRed, "Client", "Eroare conectare: "+err.Error())
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(task + "\n"))
	if err != nil {
		logLine(colorRed, "Client", "Eroare trimitere task: "+err.Error())
		return
	}

	reader := bufio.NewReader(conn)
	result, err := reader.ReadString('\n')
	if err != nil {
		logLine(colorRed, "Client", "Eroare primire răspuns: "+err.Error())
		return
	}

	logLine(colorGreen, "Client", fmt.Sprintf("Sent: %-15s | Received: %s", task, strings.TrimSpace(result)))
}

func main() {
	var wg sync.WaitGroup
	tasks := generateTasks()

	for _, task := range tasks {
		wg.Add(1)
		go sendTask(task, &wg)
		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()
}
