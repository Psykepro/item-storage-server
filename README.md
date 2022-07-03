# item-storage-server
The server for the item-storage project. It uses RabbitMQ Queue to receive requests in format of protobuf messages. 
And logs responses as protobuf messages marshaled/serialized as bytes in `responses.txt` which is located at project level.

## Requirements

#### Install Go
Instructions: https://go.dev/doc/install

#### Install and Start RabbitMQ (if not using docker-compose):
Instructions: https://www.rabbitmq.com/download.html

## Run with docker:
### 1. Locally:
#### 1.1 Start RabbitMQ from "docker-compose.local.yaml":`
`docker-compose up --build -f ./docker-compose.local.yaml`

#### 1.2. In `./config/config.yml` change property `rabbitmq.host` to `localhost`:
#### 1.3. Run the app:
`go run cmd/main.go`

### 1. Run with docker the app and the RabbitMQ:
#### 1.1 Start RabbitMQ from "docker-compose.yaml":`
`docker-compose up --build -f ./docker-compose.yaml`

