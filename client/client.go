package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	host        = flag.String("host", "127.0.0.1", "Server host")
	port        = flag.Int("port", 50051, "Server port")
	numTasks    = flag.Int("tasks", 10, "Number of tasks to generate")
	delayMillis = flag.Int("delay", 100, "Delay between task submissions (ms)")
	timeout     = flag.Int("timeout", 5, "Response timeout (seconds)")
	logFile     = flag.String("log", "app-logs/client.log", "Path to client log file")
)

func init() {
	f, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))
	log.SetFlags(log.Ldate | log.Ltime)
	rand.Seed(time.Now().UnixNano())
}

func generateTasks(n int) []string {
	ops := []string{"Add", "Minus", "Multiply", "Subtract", "Power", "Mod"}
	prios := []string{"LOW", "MEDIUM", "HIGH"}
	var tasks []string
	for i := 0; i < n; i++ {
		a, b := rand.Intn(100), rand.Intn(100)
		t := fmt.Sprintf("%s %d %d %s", ops[i%len(ops)], a, b, prios[i%len(prios)])
		tasks = append(tasks, t)
	}
	return tasks
}

func sendTask(task string, wg *sync.WaitGroup) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", *host, *port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	log.Println("Sending task:", task)
	conn.Write([]byte(task + "\n"))

	conn.SetReadDeadline(time.Now().Add(time.Duration(*timeout) * time.Second))
	resp, _ := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("Read error (timeout=%ds): %v", *timeout, err)
	} else {
		log.Println("Received:", strings.TrimSpace(resp))
	}
	conn.Write([]byte("UNREGISTER_CLIENT\n"))
}

func main() {
	flag.Parse()
	tasks := generateTasks(*numTasks)
	var wg sync.WaitGroup
	for _, t := range tasks {
		wg.Add(1)
		go sendTask(t, &wg)
		time.Sleep(time.Duration(*delayMillis) * time.Millisecond)
	}
	wg.Wait()
	log.Println("All tasks completed.")
}