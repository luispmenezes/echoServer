FROM golang:latest as build-stage
WORKDIR /echoServer
COPY ./go.mod ./
RUN go mod download
COPY ./main.go ./main.go
RUN go build

FROM ubuntu:latest as production-stage
RUN apt update &&\
    apt install -y ca-certificates &&\
    update-ca-certificates 2>/dev/null || true
WORKDIR /echoServer
COPY --from=build-stage /echoServer/echoServer /echoServer/echoServer
ENTRYPOINT ["/echoServer/echoServer"]