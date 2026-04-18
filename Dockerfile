FROM golang:alpine AS builder
WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd/server


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 7860

CMD ["./main"]
