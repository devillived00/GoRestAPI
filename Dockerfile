FROM golang:latest

RUN  mkdir /build
WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./


RUN export GO111MODULE=on
RUN go get github.com/devillived00/GoRestAPI/main
RUN cd /build && git clone https://github.com/devillived00/GoRestAPI.git

RUN cd /build/GoRestAPI/main && go build

EXPOSE 8000

ENTRYPOINT [ "/build/GoRestAPI/main/restAPI" ]