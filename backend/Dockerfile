FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/sorting ./cmd/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/sorting /bin/sorting
COPY migrations /migrations

EXPOSE 8080

CMD ["/bin/sorting"]
