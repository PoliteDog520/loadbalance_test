#======================================#
# builder env
#

FROM golang:1.15 AS builder

WORKDIR /go/src/app
COPY .  .

RUN go mod tidy
RUN go env -w GO111MODULE="on"
RUN GOOS=linux   GOARCH=amd64  CGO_ENABLED=0 \
    go build -mod=mod -o server  server.go
#     go build -mod vendor -v -o server server.go


#======================================#
# deployment env
#

FROM alpine:3.11.3

# copy executable
COPY --from=builder /go/src/app/server  /

EXPOSE 80
EXPOSE 30051

# Define default command.
ENTRYPOINT ["./server"]