
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
