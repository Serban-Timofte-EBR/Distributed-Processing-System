package com.distributed.server;

import com.distributed.generated.WorkerRegisterRequest;
import com.distributed.generated.WorkerRegisterResponse;
import com.distributed.generated.WorkerServiceGrpc;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;
import java.io.IOException;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.logging.Level;
import java.util.logging.Logger;

public class TaskServer {
    private static final Logger logger = Logger.getLogger(TaskServer.class.getName());
    private final Server server;
    private final int port;
    private final Map<String, Integer> registeredWorkers = Collections.synchronizedMap(new HashMap<>()); // Worker ID -> Port

    public TaskServer(int port) {
        this.port = port;
        this.server = ServerBuilder.forPort(port)
                .addService(new WorkerRegistrationService(this))
                .build();
    }

    public void registerWorker(String workerId, int workerPort) {
        registeredWorkers.put(workerId, workerPort);
        logger.info("Worker registered: " + workerId + " on port " + workerPort + " (Total Workers: " + registeredWorkers.size() + ")");
    }

    public void unregisterWorker(String workerId) {
        registeredWorkers.remove(workerId);
        logger.info("Worker unregistered: " + workerId + " (Remaining Workers: " + registeredWorkers.size() + ")");
    }

    public Integer getNextAvailableWorker() {
        if (registeredWorkers.isEmpty()) {
            logger.warning("No available workers!");
            return null;
        }
        return registeredWorkers.values().iterator().next(); // Simple round-robin logic
    }

    public void start() throws IOException {
        server.start();
        logger.info("Task Server started on port " + port + "...");
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            logger.warning("Shutting down gRPC server...");
            TaskServer.this.stop();
        }));
    }

    public void stop() {
        if (server != null) {
            server.shutdown();
            logger.info("gRPC Server stopped.");
        }
    }

    public void awaitTermination() throws InterruptedException {
        if (server != null) {
            server.awaitTermination();
        }
    }

    public static void main(String[] args) throws IOException, InterruptedException {
        int port = 50051;
        TaskServer server = new TaskServer(port);
        server.start();
        server.awaitTermination();
    }
}
