FROM golang:1.23-alpine

RUN apk add --no-cache git git-lfs sqlite strace

WORKDIR /var/opt/tester
COPY . /var/opt/tester
RUN git lfs pull

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
RUN chmod +x main

WORKDIR /app
CMD ["strace", "/var/opt/tester/main"]

# CMD ["/var/opt/tester/main"]

# RUN rm -f "/app/test-1.db" || true
# # CMD ["strace", "ln", "/var/opt/tester/companies.db", "/app/test-1.db"]

# RUN ln -s /var/opt/tester/companies.db /app/test-1.db
# CMD ["strace", "sqlite3", "/app/test-1.db", ".eqp full", "SELECT id, name FROM companies LIMIT 1"]