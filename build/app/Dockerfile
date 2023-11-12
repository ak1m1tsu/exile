FROM golang:alpine AS builder
RUN apk add alpine-sdk
RUN apk add --no-cache git
WORKDIR /src/service
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o /bin/service ./cmd/app/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/service /service
EXPOSE ${PORT}
ENTRYPOINT /service
