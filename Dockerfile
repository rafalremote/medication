FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/server ./cmd/server

ENV PORT=8080

EXPOSE ${PORT}

ENTRYPOINT ["./server"]
