FROM golang:1.12.7 as builder

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -o philote .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/philote .
ENTRYPOINT ["./philote"]
