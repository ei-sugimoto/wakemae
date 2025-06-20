FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o wakemae .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/wakemae .

CMD ["./wakemae", "serve"]