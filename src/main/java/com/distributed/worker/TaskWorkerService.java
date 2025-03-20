package com.distributed.worker;

import com.distributed.generated.TaskServiceGrpc;
import com.distributed.generated.TaskRequest;
import com.distributed.generated.TaskResponse;
import io.grpc.stub.StreamObserver;
import java.util.logging.Level;
import java.util.logging.Logger;

public class TaskWorkerService extends TaskServiceGrpc.TaskServiceImplBase {
    private static final Logger logger = Logger.getLogger(TaskWorkerService.class.getName());
    private final String workerId;

    public TaskWorkerService(String workerId) {
        this.workerId = workerId;
    }

    @Override
    public void executeTask(TaskRequest request, StreamObserver<TaskResponse> responseObserver) {
        String taskName = request.getTaskName();
        logger.info("Worker " + workerId + " executing task: " + taskName);

        // Simulated execution
        String result = "Executed task: " + taskName;
        int exitCode = 0;

        TaskResponse response = TaskResponse.newBuilder()
                .setExitCode(exitCode)
                .setResult(result)
                .build();

        responseObserver.onNext(response);
        responseObserver.onCompleted();
        logger.info("Worker " + workerId + " completed task: " + result);
    }
}
