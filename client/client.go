package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var tasks = []string{
	"Multiply 3 7",
	"Multiply 2 8",
	"Multiply 4 5",
	"Multiply 10 6",
	"Multiply 9 2",
}

// ✅ Sends a task to the server
func sendTask(task string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "127.0.0.1:50051")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(task + "\n"))
	if err != nil {
		fmt.Println("Error sending task:", err)
		return
	}

	// ✅ Wait for response from server
	reader := bufio.NewReader(conn)
	result, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error receiving response:", err)
		return
	}

	fmt.Printf("[Client] Sent: %s | Received: %s", task, strings.TrimSpace(result))
}

func main() {
	var wg sync.WaitGroup

	for _, task := range tasks {
		wg.Add(1)
		go sendTask(task, &wg)
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
}
