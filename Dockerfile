# syntax=docker/dockerfile:1

FROM golang:1.16

RUN  mkdir /build
WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN export GO111MODULE=on

RUN cd /build && git clone https://github.com/devillived00/GoRestAPI.git

RUN cd /build/GoRestAPI/main && go build restAPI.go

EXPOSE 8080

CMD [ "/build/GoRestAPI/main/restAPI" ]
