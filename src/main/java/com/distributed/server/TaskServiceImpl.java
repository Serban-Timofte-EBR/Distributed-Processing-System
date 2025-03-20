package com.distributed.server;

import com.distributed.generated.TaskServiceGrpc;
import com.distributed.generated.TaskRequest;
import com.distributed.generated.TaskResponse;
import io.grpc.stub.StreamObserver;
import java.util.logging.Level;
import java.util.logging.Logger;

public class TaskServiceImpl extends TaskServiceGrpc.TaskServiceImplBase {
    private static final Logger logger = Logger.getLogger(TaskServiceImpl.class.getName());
    private final TaskServer server;

    public TaskServiceImpl(TaskServer server) {
        this.server = server;
    }

    @Override
    public void executeTask(TaskRequest request, StreamObserver<TaskResponse> responseObserver) {
        String workerPort = server.getNextAvailableWorker();
        if (workerPort == null) {
            logger.warning("No workers available to execute the task.");
            responseObserver.onError(new RuntimeException("No workers available"));
            return;
        }

        logger.info("Forwarding task to worker at port: " + workerPort);

        // In a real-world scenario, the server would send this request to the worker via another gRPC call
        // Simulating the worker's execution
        String taskName = request.getTaskName();
        String result = "Executed task: " + taskName;
        int exitCode = 0; // Simulating successful execution

        TaskResponse response = TaskResponse.newBuilder()
                .setExitCode(exitCode)
                .setResult(result)
                .build();

        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }
}
