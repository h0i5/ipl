FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .

FROM alpine:latest
RUN apk add --no-cache openssh-keygen
WORKDIR /app
COPY --from=builder /app/server .
RUN mkdir -p .ssh
EXPOSE 22
EXPOSE 6767
CMD ["./server"]
