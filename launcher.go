package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

const (
	numClients = 5
	clientPath = "client/client.go"
	logDir     = "logs"
)

func runClient(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	logFilePath := fmt.Sprintf("%s/Client-%d.log", logDir, id)
	logFile, err := os.Create(logFilePath)
	if err != nil {
		fmt.Printf("[Launcher] Error creating log file for Client-%d: %v\n", id, err)
		return
	}
	defer logFile.Close()

	cmd := exec.Command("go", "run", clientPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("[Launcher] Error getting stdout for Client-%d: %v\n", id, err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("[Launcher] Error getting stderr for Client-%d: %v\n", id, err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("[Launcher] Failed to start Client-%d: %v\n", id, err)
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("[Client-%d] %s\n", id, line)
			logFile.WriteString(line + "\n")
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("[Client-%d][ERROR] %s\n", id, line)
			logFile.WriteString("[ERROR] " + line + "\n")
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Printf("[Launcher] Client-%d exited with error: %v\n", id, err)
	}
}

func main() {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			fmt.Println("[Launcher] Failed to create logs directory:", err)
			return
		}
	}

	var wg sync.WaitGroup
	fmt.Printf("[Launcher] Launching %d parallel clients...\n\n", numClients)

	for i := 1; i <= numClients; i++ {
		wg.Add(1)
		go runClient(i, &wg)
	}

	wg.Wait()
	fmt.Println("\n[Launcher] All clients finished.")
}
