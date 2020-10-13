FROM golang:1.14-alpine as builder
RUN apk update && apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o bin/bca-rate cmd/bca-rate/*.go

FROM alpine:3
RUN apk add --no-cache ca-certificates
COPY --from=builder /src/bin bin
USER nobody:nobody
ENTRYPOINT [ "/bin/bca-rate" ]
