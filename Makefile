
CURPATH=$(PWD)
export GOBIN=$(CURDIR)/bin
export PATH:=$(GOBIN):$(PATH)

IMAGE_TAG?=huikang/go-es-ocp:latest
OPERATOR_NAMESPACE=openshift-operators-redhat

all: build

fmt:
	@gofmt -s -l -w *.go

build: fmt
	@go build -o $(GOBIN)/go-es-ocp ./main.go

image:
	@if [ $${SKIP_BUILD:-false} = false ] ; then \
		docker build -t $(IMAGE_TAG) . ; \
	fi

deploy-image: image
	docker push $(IMAGE_TAG)

deploy:
	echo "Deploy to elasticsearch operator namespace"
	kubectl -n $(OPERATOR_NAMESPACE) apply -f ./k8s.yaml

undeploy:
	kubectl -n $(OPERATOR_NAMESPACE) delete -f ./k8s.yaml
