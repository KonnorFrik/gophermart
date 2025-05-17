FROM golang:alpine

EXPOSE 8000

COPY ./user_service /go/src/
COPY container/init.sh /go/src/
WORKDIR /go/src
RUN sh init.sh
WORKDIR /go
USER user

CMD [ "bin/gophermart" ]

