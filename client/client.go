package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func generateTasks() []string {
	var tasks []string
	for i := 1; i <= 50; i++ {
		a := i
		b := i + 1
		tasks = append(tasks, fmt.Sprintf("Multiply %d %d", a, b))
	}
	return tasks
}

func sendTask(task string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "127.0.0.1:50051")
	if err != nil {
		fmt.Println("Eroare conectare:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(task + "\n"))
	if err != nil {
		fmt.Println("Eroare trimitere task:", err)
		return
	}

	reader := bufio.NewReader(conn)
	result, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Eroare primire răspuns:", err)
		return
	}

	fmt.Printf("\n[Client] Sent: %-15s | Received: %s", task, strings.TrimSpace(result))
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
