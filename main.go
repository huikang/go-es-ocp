package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	elasticsearch6 "github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/estransport"
	"github.com/tidwall/gjson"
)

func main() {
	fmt.Println("Elasticserach go client testing...")

	var es_address = flag.String("es_addr", "http://localhost:9200",
		"elasticsearch address (default: localhost:9200)")
	flag.Parse()

	if *es_address == "" {
		log.Fatalf("es address is empty")
	}
	log.Printf("es address: %s\n", *es_address)

	// es, err := elasticsearch6.NewDefaultClient()
	cfg := elasticsearch6.Config{
		Addresses: []string{*es_address},
	}
	es, err := elasticsearch6.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s\n", err)
	}

	// Get cluster info
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting the response: %s\n", err)
	}
	defer res.Body.Close()
	log.Print(es.Transport.(*estransport.Client).URLs())

	// Get cluster version
	res, err = es.Cluster.Stats(es.Cluster.Stats.WithPretty())
	if err != nil {
		log.Fatalf("Error getting the cluster response: %s\n", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	str := string(body)
	version := gjson.Get(str, "nodes.versions")
	log.Println(version)
}
