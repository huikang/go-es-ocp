
CURPATH=$(PWD)
export GOBIN=$(CURDIR)/bin
export PATH:=$(GOBIN):$(PATH)

IMAGE_TAG?=huikang/go-es-ocp:latest
OPERATOR_NAMESPACE=openshift-operators-redhat
ES_CONTAINER_NAME=elasticsearch
ES_OCP_IMAGE_TAG=quay.io/openshift/origin-logging-elasticsearch6
ES_IMAGE_TAG=elasticsearch:6.8.12
ES_OD_IMAGE_TAG=amazon/opendistro-for-elasticsearch:0.10.0

ES_ADDR?=http://localhost:9200
ES_TOKEN?=token

all: build

fmt:
	@gofmt -s -l -w *.go

build: fmt
	@go build -o $(GOBIN)/go-es-ocp ./main.go

run-local:
	@go run ./main.go -es_addr=$(ES_ADDR)

run-local-ocp:
	@go run ./main.go -es_addr=$(ES_ADDR) -t $(ES_TOKEN) \
		-r ./admin-ca -c ./admin-cert -k ./admin-key

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
		$(ES_IMAGE_TAG)

# -v $(PWD)/secret:/etc/elasticsearch/secret 
run-es-ocp:
	docker run -d --name $(ES_CONTAINER_NAME) \
		-v $(PWD)/etc-elasticsearch:/etc/elasticsearch \
		-v $(PWD)/config:/usr/share/java/elasticsearch/config \
		-v $(PWD)/persistent:/elasticsearch/persistent \
		-p 9200:9200 -p 9300:9300 \
		-e "CLUSTER_NAME=local" \
		-e "IS_MASTER=true" \
		-e "DC_NAME=es" \
		-e "HAS_DATA=true" \
		-e "HEAP_DUMP_LOCATION=/elasticsearch/persistent/heapdump.hprof" \
		$(ES_OCP_IMAGE_TAG)

run-es-od:
	docker run -d --name $(ES_CONTAINER_NAME) \
		-v $(PWD)/etc-elasticsearch:/etc/elasticsearch \
		-v $(PWD)/config:/usr/share/java/elasticsearch/config \
		-v $(PWD)/persistent:/elasticsearch/persistent \
		-p 9200:9200 -p 9300:9300 \
		-e "discovery.type=single-node" \
		$(ES_OD_IMAGE_TAG)

clean:
	docker rm -f $(ES_CONTAINER_NAME)
