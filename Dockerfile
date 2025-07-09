FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dialog_service ./cmd


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/dialog_service .

EXPOSE 9001


ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_NAME=qwerty
ENV GRPC_PORT=9001

CMD ["./dialog_service"]