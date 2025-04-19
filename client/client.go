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
	priorities := []string{"LOW", "MEDIUM", "HIGH"}
	operations := []string{"Add", "Minus", "Multiply", "Subtract"}

	for i := 1; i <= 10; i++ {
		a := i
		b := i + 1
		priority := priorities[i%3]  // switch between LOW, MEDIUM, HIGH
		operation := operations[i%4] // cycle through Add, Minus, Multiply, Subtract
		tasks = append(tasks, fmt.Sprintf("%s %d %d %s", operation, a, b, priority))
	}
	return tasks
}

func sendTask(task string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "127.0.0.1:50051")
	if err != nil {
		fmt.Println("[Client] Eroare conectare:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(task + "\n"))
	if err != nil {
		fmt.Println("[Client] Eroare trimitere task:", err)
		return
	}

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("[Client] Failed to read response:", err)
		return
	}
	fmt.Println("[Client] Received:", strings.TrimSpace(response))

	fmt.Fprintf(conn, "UNREGISTER_CLIENT\n")
}

func main() {
	var wg sync.WaitGroup
	tasks := generateTasks()
	conn, err := net.Dial("tcp", "127.0.0.1:50051")
	if err != nil {
		fmt.Println("Client: Nu s-a putut conecta la server:", err)
		return
	}

	defer func() {
		conn.Write([]byte("UNREGISTER_CLIENT\n")) // ✅ Trimite semnalul explicit
		conn.Close()
		fmt.Println("Client: Connection finished.")
	}()

	for _, task := range tasks {
		wg.Add(1)
		go sendTask(task, &wg)
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
}
