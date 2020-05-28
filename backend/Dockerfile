FROM golang:1.14 as builder
WORKDIR $GOPATH/src/subscription-website/backend
COPY ./ .

RUN GOOS=linux GOARCH=386 go get .
RUN GOOS=linux GOARCH=386 go build -ldflags="-w -s" -v

RUN cp backend /

FROM alpine:latest
COPY --from=builder /backend /
COPY  db/migrations/ /db/migrations/

CMD ["/backend"]
