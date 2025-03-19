# **Distributed Processing System**

## **Project Overview**
This project implements a **distributed task execution system** using a **client-server architecture**. Clients can both **submit and process tasks**, making it a **fully decentralized system** with dynamic workload distribution and efficient resource utilization.

## **System Components**
1. **Server**
   - Listens on **loopback (`127.0.0.1`)** on a fixed port (e.g., `5000`).
   - Manages a list of **active processing clients**.
   - Receives **task execution requests** from clients.
   - Assigns tasks to clients using a **load balancing algorithm**.
   - Returns **execution results** to the requesting client.

2. **Client (Worker Node)**
   - **Registers** with the server on startup, sending its available processing port.
   - **Listens on its own port** for incoming processing requests.
   - **Processes assigned tasks locally** and returns the execution result.
   - **Deregisters** from the server upon shutdown.

3. **Client (Task Submitter)**
   - **Sends a task execution request** to the server.
   - **Waits for and receives the execution result** from the server.

## **Workflow**
1. **Client (Worker Node) starts and registers with the server:**
   ```
   REGISTER 127.0.0.1:<PORT>
   ```
2. **Client (Task Submitter) sends a task to the server:**
   ```
   EXECUTE {
     "task_id": "1234",
     "function": "multiply",
     "arguments": [5, 10]
   }
   ```
3. **Server selects an available client and forwards the task:**
   ```
   RUN {
     "task_id": "1234",
     "function": "multiply",
     "arguments": [5, 10]
   }
   ```
4. **Processing client executes the task and returns the result:**
   ```
   RESULT {
     "task_id": "1234",
     "exit_code": 0,
     "result": 50
   }
   ```
5. **Server forwards the result back to the requesting client:**
   ```
   TASK_COMPLETED {
     "task_id": "1234",
     "exit_code": 0,
     "result": 50
   }
   ```
6. **Client (Worker Node) deregisters when shutting down:**
   ```
   DEREGISTER 127.0.0.1:<PORT>
   ```

## **Load Balancing Strategies**
- **Round Robin** – Assigns tasks to clients in a cyclic order.
- **Least Connections** – Assigns tasks to the client with the fewest active tasks.
- **Priority Queues** – Prioritizes specific types of tasks based on importance.

## **Technologies**
- **Networking:** gRPC, WebSockets, or raw TCP/IP Sockets
- **Programming Language:** Java, C#, or Go
- **Process Management:** In-memory task execution with defined function mappings
- **Concurrency:** Threading or Async processing (e.g., `ThreadPool`, `Task`, `goroutines`)

## **System Development**
- Create the **server-client communication** using a chosen networking protocol.
- Create a **registration system** for worker nodes.
- Develop **task execution logic and load balancing**.
- Implement **result handling and reporting**.
- Add **error handling, logging, and performance optimization**.

## **Example Diagram**
```plaintext
[Client A] ---> (EXECUTE Task) ---> [Server] ---> (Assign Task) ---> [Client B]
[Client B] ---> (PROCESS Task) ---> [Server] ---> (TASK_COMPLETED) ---> [Client A]
```
