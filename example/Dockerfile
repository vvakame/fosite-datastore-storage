FROM golang:1.12 as builder

WORKDIR /go/src/github.com/vvakame/fosite-datastore-storage/example
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=off go build -v -o main


FROM alpine

RUN apk --no-cache --update add ca-certificates

COPY --from=builder /go/src/github.com/vvakame/fosite-datastore-storage/example/main /main
COPY --from=builder /go/src/github.com/vvakame/fosite-datastore-storage/example/public /public

CMD ["/main"]
