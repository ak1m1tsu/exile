#build stage
FROM golang:alpine AS builder
RUN apk add alpine-sdk
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -tags musl -o /go/bin/service ./cmd/service/main.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/service /service
ENTRYPOINT /service
