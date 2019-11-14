FROM golang:1.11 as builder
WORKDIR $GOPATH/src/github.com/AndrewDonelson/rest-svc
COPY ./ .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -v
RUN cp rest-api /

FROM alpine:latest
COPY --from=builder /rest-api /
CMD ["/rest-api"]
