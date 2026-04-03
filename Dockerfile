FROM golang:1.24-alpine AS builder

WORKDIR /app

# install git for some go modules
RUN apk add --no-cache git

# copy go module files first for better caching
COPY go.mod ./
COPY go.sum* ./

# download dependencies
RUN go mod download

COPY . .

# build binary
RUN go build -o main ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main ./main

EXPOSE 8080

CMD ["./main"]