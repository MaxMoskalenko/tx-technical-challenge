FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/tx-technical-challenge .

FROM alpine:3.22

WORKDIR /app

ENV GRPC_ADDR=:50051
EXPOSE 50051

COPY --from=builder /bin/tx-technical-challenge /usr/local/bin/tx-technical-challenge

CMD ["tx-technical-challenge"]
