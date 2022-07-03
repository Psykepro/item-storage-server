FROM golang:1.18

#ENV TZ=Europe/Moscow
#RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV config=docker
ENV GO111MODULE=on

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go install github.com/githubnemo/CompileDaemon@latest

ENTRYPOINT CompileDaemon --build="go build cmd/main.go" --command=./main