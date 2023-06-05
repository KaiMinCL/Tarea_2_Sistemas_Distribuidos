Tarea n° 2: TrustBank

@authors: 
 - David ALVAREZ
 - Jules EVANS


To execute the project one must launch 4 terminals: one for the client, one for the server and the tow others for the rabbitmq and grpc.

First you need to enter the main folder in each terminal: `cd ./Tarea2`

Now in each terminal you need to acces the correct folder for executing the code: 
 - For the client in the folder client: ` cd ./client`
 - For the server int the api folder: ` cd ./server`
 - For the rabbitmq in the rabbitmq folder: `cd ./rabbitmq`
 - For the grpc in the grpc folder: `cd ./grpc`

Next in each respective terminal execute in order: 

 - In the client one : `go run client_menu.go`
 - In the server one: `go run server_api.go`
 - In the rabbitmq one: `go run rabbitmq_consumer.go`
 - In the grpc one: `go run server_grpc.go`

 Now the program is running, follow the instructions of the client terminal.
