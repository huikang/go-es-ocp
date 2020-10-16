package main

import (
	"fmt"

	elasticsearch6 "github.com/elastic/go-elasticsearch/v6"
)

func main() {
	fmt.Println("Elasticserach go client testing...")

	es6, _ := elasticsearch6.NewDefaultClient()
	fmt.Printf("%v", es6.Cluster)
}
