
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

0. Install the CLO operator and ES operator, and create an CLO instance.

1. Expose the log store

- Make sure you have `oc login` to the cluster

```bash
    ./expose-log-store.sh

    # Verify access:
    token=$(oc whoami -t)
    routeES=`oc get route elasticsearch -o jsonpath={.spec.host}`
    curl -tlsv1.2 --insecure -H "Authorization: Bearer ${token}" "https://${routeES}/.operations.*/_search?size=1" | jq
```

2. Run test program:

```bash
    ES_ADDR=https://${routeES}:443 TOKEN=<sa token> make run-local-ocp
```
