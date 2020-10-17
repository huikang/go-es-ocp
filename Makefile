
CURPATH=$(PWD)
export GOBIN=$(CURDIR)/bin
export PATH:=$(GOBIN):$(PATH)

IMAGE_TAG?=huikang/go-es-ocp:latest
OPERATOR_NAMESPACE=openshift-operators-redhat
ES_CONTAINER_NAME=elasticsearch

all: build

fmt:
	@gofmt -s -l -w *.go

build: fmt
	@go build -o $(GOBIN)/go-es-ocp ./main.go

run-local:
	@go run ./main.go -es_addr="http://localhost:9200"

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

run-es:
	docker run -d --name $(ES_CONTAINER_NAME) \
		-p 9200:9200 -p 9300:9300 \
		-e "discovery.type=single-node" \
		elasticsearch:6.8.12

clean:
	docker rm -f $(ES_CONTAINER_NAME)
