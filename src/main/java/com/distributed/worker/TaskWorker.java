package com.distributed.worker;

import com.distributed.generated.WorkerServiceGrpc;
import com.distributed.generated.WorkerRegisterRequest;
import com.distributed.generated.WorkerRegisterResponse;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import java.io.IOException;
import java.util.UUID;
import java.util.logging.Level;
import java.util.logging.Logger;

public class TaskWorker {
    private static final Logger logger = Logger.getLogger(TaskWorker.class.getName());
    private final Server server;
    private final String workerId;
    private final int workerPort;
    private final ManagedChannel channel;
    private final WorkerServiceGrpc.WorkerServiceBlockingStub serverStub;

    public TaskWorker(int workerPort, String serverHost, int serverPort) {
        this.workerId = UUID.randomUUID().toString();
        this.workerPort = workerPort;
        this.server = ServerBuilder.forPort(workerPort)
                .addService(new TaskWorkerService(workerId))
                .build();

        this.channel = ManagedChannelBuilder.forAddress(serverHost, serverPort)
                .usePlaintext()
                .build();
        this.serverStub = WorkerServiceGrpc.newBlockingStub(channel);
    }

    public void registerWithServer() {
        logger.info("Registering worker with server...");
        WorkerRegisterRequest request = WorkerRegisterRequest.newBuilder()
                .setWorkerId(workerId)
                .setWorkerPort(workerPort)
                .build();

        WorkerRegisterResponse response = serverStub.registerWorker(request);
        logger.info("Worker registered successfully with server. Assigned ID: " + response.getConfirmationMessage());
    }

    public void start() throws IOException {
        server.start();
        logger.info("Worker started with ID: " + workerId + " on port " + workerPort);
        registerWithServer();
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            logger.warning("Shutting down worker: " + workerId);
            stop();
        }));
    }

    public void stop() {
        if (server != null) {
            server.shutdown();
            logger.info("Worker " + workerId + " stopped.");
        }
    }

    public void awaitTermination() throws InterruptedException {
        if (server != null) {
            server.awaitTermination();
        }
    }

    public static void main(String[] args) throws IOException, InterruptedException {
        if (args.length < 2) {
            System.err.println("Usage: TaskWorker <workerPort> <serverHost> <serverPort>");
            System.exit(1);
        }

        int workerPort = Integer.parseInt(args[0]);
        String serverHost = args[1];
        int serverPort = Integer.parseInt(args[2]);

        TaskWorker worker = new TaskWorker(workerPort, serverHost, serverPort);
        worker.start();
        worker.awaitTermination();
    }
}
