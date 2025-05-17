# === Build stage === #
FROM golang:alpine AS builder
COPY ./user_service/go.mod ./user_service/go.sum /go/src/
WORKDIR /go/src
RUN go mod download
COPY ./user_service /go/src/
RUN go build -o ../bin/.


# === Final stage === #
FROM golang:alpine

EXPOSE 8000

COPY --from=builder /go/bin/gophermart /go/bin/gophermart
WORKDIR /go
RUN adduser -S user
USER user

CMD [ "bin/gophermart" ]

