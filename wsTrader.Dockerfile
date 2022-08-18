FROM golang:1.18

#RUN apt update && apt install -y vim net-tools htop

WORKDIR /usr/src/app

COPY . .
RUN go mod download && go mod verify
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app","-addr",":2096"]

EXPOSE 2096