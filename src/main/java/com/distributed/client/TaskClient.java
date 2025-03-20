package com.distributed.client;

import com.distributed.generated.TaskRequest;
import com.distributed.generated.TaskResponse;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.stub.StreamObserver;
import java.util.Arrays;
import java.util.logging.Level;
import java.util.logging.Logger;

public class TaskClient {
    private static final Logger logger = Logger.getLogger(TaskClient.class.getName());
    private final ManagedChannel channel;

    public TaskClient(String serverHost, int serverPort) {
        this.channel = ManagedChannelBuilder.forAddress(serverHost, serverPort)
                .usePlaintext()
                .build();
    }

    public void sendTask(String taskName, String... args) {
        logger.info("Sending task to server: " + taskName + " with arguments " + Arrays.toString(args));

        TaskRequest request = TaskRequest.newBuilder()
                .setTaskName(taskName)
                .addAllArguments(Arrays.asList(args))
                .build();

        StreamObserver<TaskResponse> responseObserver = new StreamObserver<>() {
            @Override
            public void onNext(TaskResponse response) {
                logger.info("Received response: Exit Code = " + response.getExitCode() +
                        ", Result = " + response.getResult());
            }

            @Override
            public void onError(Throwable t) {
                logger.severe("Error processing task: " + t.getMessage());
            }

            @Override
            public void onCompleted() {
                logger.info("Task processing complete.");
                channel.shutdown();
            }
        };
    }

    public static void main(String[] args) {
        TaskClient client = new TaskClient("localhost", 50051);
        client.sendTask("Multiply", "5", "10");
    }
}
