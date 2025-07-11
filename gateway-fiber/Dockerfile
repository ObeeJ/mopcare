FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY gateway-fiber/go.mod gateway-fiber/go.sum ./
RUN go mod download

COPY gateway-fiber/ .
RUN go build -o gateway-fiber .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/gateway-fiber .

EXPOSE 9090
CMD ["./gateway-fiber"]