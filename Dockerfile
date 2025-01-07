FROM golang:1.23 as builder

ARG GOPROXY
ENV GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0 \
    GO111MODULE=on

WORKDIR /go/src/github.com/heczzots/denet

ADD go.* ./
RUN go mod download

ADD . ./
WORKDIR cmd
RUN go build -v

FROM alpine:latest

WORKDIR /app
COPY --from=builder /go/src/github.com/heczzots/denet/cmd/cmd .
COPY --from=builder /go/src/github.com/heczzots/denet/migrations /app/migrations
EXPOSE 8080

ENTRYPOINT ["./cmd"]