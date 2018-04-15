FROM golang:1.9.0 as builder

WORKDIR /go/src/github.com/sbspider/shipper/consignment-service

COPY . .

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep init -v && dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY --from=builder /go/src/github.com/sbspider/shipper/consignment-service/consignment-service .

CMD ["./consignment-service"]
