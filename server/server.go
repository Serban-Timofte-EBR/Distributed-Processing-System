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
	"sync"
)

type Task struct {
	ID        string
	Operation string
	Args      []string
	Priority  string
}

var (
	host     = flag.String("host", "127.0.0.1", "Server host address")
	port     = flag.Int("port", 50051, "Server listening port")
	logFile  = flag.String("log", "app-logs/server.log", "Path to server log file")
	queueMu  sync.Mutex
	highQ    []Task
	mediumQ  []Task
	lowQ     []Task
	results  = make(map[string]net.Conn)
)

func initLogger() {
	f, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file %s: %v", *logFile, err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}

func logInfo(label, msg string) {
	log.Printf("[SERVER][%s] %s", label, msg)
}

func addTask(conn net.Conn, taskLine string) {
	parts := strings.Fields(taskLine)
	if len(parts) < 4 {
		logInfo("ERROR", "Invalid task format: "+taskLine)
		return
	}
	op := parts[0]
	args := parts[1 : len(parts)-1]
	prio := strings.ToUpper(parts[len(parts)-1])
	id := strconv.Itoa(len(results) + 1)
		task := Task{ID: id, Operation: op, Args: args, Priority: prio}

	queueMu.Lock()
	switch prio {
	case "HIGH":
		highQ = append(highQ, task)
	case "LOW":
		lowQ = append(lowQ, task)
	default:
		mediumQ = append(mediumQ, task)
	}
	results[id] = conn
	queueMu.Unlock()

	logInfo("TASK_ADDED", fmt.Sprintf("%s %v [ID:%s] Priority:%s", op, args, id, prio))
}

func getNext() *Task {
	queueMu.Lock()
	defer queueMu.Unlock()
	if len(highQ) > 0 {
		t := highQ[0]
		highQ = highQ[1:]
		return &t
	}
	if len(mediumQ) > 0 {
		t := mediumQ[0]
		mediumQ = mediumQ[1:]
		return &t
	}
	if len(lowQ) > 0 {
		t := lowQ[0]
		lowQ = lowQ[1:]
		return &t
	}
	return nil
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			logInfo("CLOSED", fmt.Sprintf("%s", err))
			return
		}
		msg := strings.TrimSpace(line)
		switch {
		case msg == "REGISTER_WORKER":
			logInfo("WORKER", "Registered " + conn.RemoteAddr().String())

		case msg == "REQUEST_TASK":
			t := getNext()
			if t != nil {
				payload := fmt.Sprintf("%s %s %s", t.ID, t.Operation, strings.Join(t.Args, " "))
				conn.Write([]byte(payload + "\n"))
				logInfo("ASSIGN", payload)
			} else {
				conn.Write([]byte("NO_TASK\n"))
			}

		case strings.HasPrefix(msg, "RESULT"):
			parts := strings.Fields(msg)
			if len(parts) < 3 {
				logInfo("ERROR", "Malformed result: "+msg)
				continue
			}
			id := parts[1]
			res := strings.Join(parts[2:], " ")
			logInfo("RESULT", fmt.Sprintf("Task %s → %s", id, res))
			if client, ok := results[id]; ok {
				client.Write([]byte(res + "\n"))
				delete(results, id)
			} else {
				logInfo("ERROR", "No client for Task ID "+id)
			}

		case msg == "UNREGISTER_CLIENT":
			logInfo("CLIENT", "Disconnected "+conn.RemoteAddr().String())
			return

		default:
			addTask(conn, msg)
		}
	}
}

func main() {
	flag.Parse()
	initLogger()
	addr := fmt.Sprintf("%s:%d", *host, *port)
	logInfo("START", "Listening on " + addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to bind to %s: %v", addr, err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			logInfo("ERROR", err.Error())
			continue
		}
		go handleConn(conn)
	}
}