package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:50051")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Enter task (e.g., 'Multiply 5 10'):")
	reader := bufio.NewReader(os.Stdin)
	task, _ := reader.ReadString('\n')
	task = strings.TrimSpace(task)

	_, err = conn.Write([]byte(task + "\n"))
	if err != nil {
		fmt.Println("Error sending task:", err)
		return
	}

	serverReader := bufio.NewReader(conn)
	result, err := serverReader.ReadString('\n')
	if err != nil {
		fmt.Println("Error receiving response:", err)
		return
	}

	result = strings.TrimSpace(result)
	fmt.Println("Received result from server:", result)
}
