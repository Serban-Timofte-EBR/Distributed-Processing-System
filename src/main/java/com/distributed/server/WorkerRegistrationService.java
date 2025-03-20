package com.distributed.server;

import com.distributed.generated.WorkerRegisterRequest;
import com.distributed.generated.WorkerRegisterResponse;
import com.distributed.generated.WorkerServiceGrpc;
import io.grpc.stub.StreamObserver;
import java.util.logging.Level;
import java.util.logging.Logger;

public class WorkerRegistrationService extends WorkerServiceGrpc.WorkerServiceImplBase {
    private static final Logger logger = Logger.getLogger(WorkerRegistrationService.class.getName());
    private final TaskServer server;

    public WorkerRegistrationService(TaskServer server) {
        this.server = server;
    }

    @Override
    public void registerWorker(WorkerRegisterRequest request, StreamObserver<WorkerRegisterResponse> responseObserver) {
        String workerId = request.getWorkerId();
        int workerPort = request.getWorkerPort();

        server.registerWorker(workerId, workerPort);
        logger.info("Registered new worker: " + workerId + " at port " + workerPort);

        WorkerRegisterResponse response = WorkerRegisterResponse.newBuilder()
                .setConfirmationMessage("Worker " + workerId + " registered successfully")
                .build();

        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }
}
