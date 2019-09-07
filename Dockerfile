FROM golang:1.12 AS builder

WORKDIR /go/src/app
COPY main.go /go/src/app/

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM scratch

COPY --from=builder /go/src/app/main main

CMD ["./main"]
