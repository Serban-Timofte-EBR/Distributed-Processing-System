package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var (
	clientPath = flag.String("client", "client/client.go", "Path to the client application")
	logDir     = flag.String("logDir", "logs", "Directory for client logs")
	numClients = flag.Int("clients", 5, "Number of clients to launch in parallel")
)

func runClient(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	logFilePath := filepath.Join(*logDir, fmt.Sprintf("Client-%d.log", id))
	logFile, err := os.Create(logFilePath)
	if err != nil {
		fmt.Printf("[Launcher] Error creating log file for Client-%d: %v\n", id, err)
		return
	}
	defer logFile.Close()

	cmd := exec.Command("go", "run", *clientPath)

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
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				extra := string(buf[:n])
				fmt.Printf("[Client-%d] %s", id, extra)
				logFile.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
	}()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				extra := string(buf[:n])
				fmt.Printf("[Client-%d][ERROR] %s", id, extra)
				logFile.WriteString("[ERROR] " + extra)
			}
			if err != nil {
				break
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Printf("[Launcher] Client-%d exited with error: %v\n", id, err)
	}
}

func main() {
	flag.Parse()

	if _, err := os.Stat(*logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(*logDir, 0755); err != nil {
			fmt.Println("[Launcher] Failed to create logs directory:", err)
			return
		}
	}

	fmt.Printf("[Launcher] Launching %d parallel clients (path: %s)...\n\n", *numClients, *clientPath)

	var wg sync.WaitGroup
	for i := 1; i <= *numClients; i++ {
		wg.Add(1)
		go runClient(i, &wg)
	}

	wg.Wait()
	fmt.Println("\n[Launcher] All clients finished.")
}