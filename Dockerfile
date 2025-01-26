# Build stage
FROM golang:1.23 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM scratch

WORKDIR /app
COPY user.csv .
COPY --from=builder /app/main .

CMD ["./main"]
