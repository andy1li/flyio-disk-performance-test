FROM golang:1.23-alpine

RUN apk add --no-cache git git-lfs

WORKDIR /var/opt/tester
COPY . /var/opt/tester
RUN git lfs pull

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
RUN chmod +x main

WORKDIR /app
CMD ["/var/opt/tester/main"]