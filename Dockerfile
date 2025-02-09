FROM golang:1.23

WORKDIR /app 

COPY . .

RUN go mod download && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main

ENTRYPOINT ["./main"]

CMD ["pinger"]
