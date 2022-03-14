FROM golang:1.16 as builder
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o dobby .

FROM alpine:3.15.0
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/dobby ./
COPY ci/config.yml .

ENTRYPOINT ["./dobby"]
