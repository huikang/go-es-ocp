FROM golang:1.15 as builder

WORKDIR /go/src/github.com/huikang/go-es-ocp
COPY . .
RUN go mod download

RUN make build 
RUN cp ./bin/go-es-ocp /usr/local/bin/


WORKDIR /usr/local/bin/
ENTRYPOINT ["go-es-ocp"]
