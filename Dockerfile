FROM golang:1.23-alpine

WORKDIR /app

# Copy source code and database
RUN mkdir -p /var/opt/tester
COPY companies.db /var/opt/tester
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
COPY main.go /app/main.go

# Build the application inside the container
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Make the binary executable
RUN chmod +x main

# Run the application
CMD ["./main"]