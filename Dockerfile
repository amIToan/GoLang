# FROM golang:1.21.6-alpine3.19
# WORKDIR /app
# COPY . .
# RUN go build -o main main.go
# EXPOSE 8080
# CMD [ "app/main" ]

# Build stage 
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
# RUN apk add curl
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY  start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8080 8000
CMD ["./main"]
ENTRYPOINT [ "./start.sh" ]
