
CURPATH=$(PWD)
export GOBIN=$(CURDIR)/bin
export PATH:=$(GOBIN):$(PATH)

all: build

fmt:
	@gofmt -s -l -w *.go

build: fmt
	@go build -o $(GOBIN)/go-es-ocp ./main.go
