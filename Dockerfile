FROM golang:1.21-alpine

WORKDIR /app

# Copy source code and database
COPY companies.db /app/companies.db
COPY go.mod /app/go.mod
COPY main.go /app/main.go

# Build the application inside the container
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Make the binary executable
RUN chmod +x main

# Run the application
CMD ["./main"]