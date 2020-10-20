# Example of how to use the es6 client in ES operator

## Run local test

1. Start elasticsearch container

```bash
    make run-es

    # To verify elasticsearch is running:
    curl localhost:9200/_cluster/stats | jq .
```

2. Run the `go-es` program. If success, the program returns the basic info of 
the cluster.

```bash
    make run-local
```

## Deploy and run against the elasticsearch in OCP

0. Install the CLO operator and ES operator, and create a CLO instance.

1. Expose the log store

- Make sure you have `oc login` to the cluster

```bash
    ./expose-log-store.sh

    token=$(oc whoami -t)
    # Expose the elasticsearch through port forwarding:
    oc -n openshift-logging port-forward service/elasticsearch 9200:9200
```

To verify the connection, open another terminal
```
    routeES=localhost:9200
    curl -tlsv1.2 --insecure -H "Authorization: Bearer ${token}" "https://${routeES}/.operations.*/_search?size=1" | jq
```

2. Run test program:

```bash
    routeES=localhost:9200
    token=<SA token from the es operator, `/var/run/secrets/kubernetes.io/serviceaccount/token`>
    ES_ADDR=https://${routeES} ES_TOKEN=${token} make run-local-ocp
```
