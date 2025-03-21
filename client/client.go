package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
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

	taskRequest := "Multiply 5 10\n"
	fmt.Println("Sent Task:", strings.TrimSpace(taskRequest))
	_, err = conn.Write([]byte(taskRequest))
	if err != nil {
		fmt.Println("Error sending task:", err)
		return
	}

	response, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Response:", strings.TrimSpace(response))
}
