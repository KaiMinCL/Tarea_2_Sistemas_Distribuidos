Tarea n° 2: TrustBank

@authors: 
 - David ALVAREZ
 - Jules EVANS

To execute the project one must launch 3 machines: One for the HTTP Server, and two more for the RabbitMq and gRPC servers. Each server can be started by using the following commands in root:
 - For the server int the api folder: `go run Tarea2/api/server_api.go`
 - For the rabbitmq in the rabbitmq folder: `go run Tarea2/rabbitmq/rabbitmq_consumer.go`
 - For the grpc in the grpc folder: `go run Tarea2/grpc/server_grpc.go`

For the client you only need to clone the `Tarea2/client` folder in the respective client and run the client with `go run Tarea2/grpc/server_grpc.go` from the root.

For the evaluation of the assigment the servers are inside three VM with the following addresses:
 - HTTP Server: `10.10.28.233:8080`
 - RabbitMq Server: `10.10.11.49:5672`
 - gRPC Server: `10.10.11.216:50051`

Because this are local network addresses, the client must be running inside the same local network of the Informatics Department.

The addresses of the servers can be modify with the `.env` file to run in any network (if the machines are setup properly).
