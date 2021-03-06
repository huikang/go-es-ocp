package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"strings"
	"time"

	elasticsearch6 "github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/estransport"
	"github.com/tidwall/gjson"
)

func getRootCA() *x509.CertPool {
	certPool := x509.NewCertPool()

	// load cert into []byte
	f := path.Join("./", "admin-ca")
	caPem, err := ioutil.ReadFile(f)
	if err != nil {
		log.Panicf("Unable to read file to get contents %v", err)
		return nil
	}
	log.Printf("ca pem %v", string(caPem))
	certPool.AppendCertsFromPEM(caPem)

	return certPool
}

func getClientCertificates() []tls.Certificate {
	certificate, err := tls.LoadX509KeyPair(
		path.Join("./", "admin-cert"),
		path.Join("./", "admin-key"),
	)
	if err != nil {
		log.Println("erro load key pairs")
		return []tls.Certificate{}
	}
	return []tls.Certificate{
		certificate,
	}
}

func oauthEsClient(esAddr, token, caPath, certPath, keyPath string) (*elasticsearch6.Client, error) {
	es := &elasticsearch6.Client{}
	httpTranport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            getRootCA(),
			// Certificates:       getClientCertificates(),
		},
	}

	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	cfg := elasticsearch6.Config{
		Header:    header,
		Addresses: []string{esAddr},
		Transport: httpTranport,
	}
	es, err := elasticsearch6.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error creating the client: %v", err)
	}
	return es, nil
}

func mTlsEsClient(esAddr, caPath, certPath, keyPath string) (*elasticsearch6.Client, error) {
	es := &elasticsearch6.Client{}
	httpTranport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            getRootCA(),
			Certificates:       getClientCertificates(),
		},
	}

	cfg := elasticsearch6.Config{
		Addresses: []string{esAddr},
		Transport: httpTranport,
	}
	es, err := elasticsearch6.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error creating the mtls client: %v", err)
	}
	return es, nil
}

func main() {
	fmt.Println("Elasticserach go client testing...")

	var esAddr = flag.String("es_addr", "http://localhost:9200",
		"elasticsearch address (default: localhost:9200)")
	var esToken = flag.String("t", "",
		"elasticsearch token")
	var caPath = flag.String("r", "",
		"CA file path")
	var certPath = flag.String("c", "",
		"Cert file path")
	var keyPath = flag.String("k", "",
		"key file path")
	flag.Parse()

	if *esAddr == "" {
		log.Fatalf("es address is empty")
	}
	log.Printf("es address: %s\n", *esAddr)

	// Setup es client
	var esClis []*elasticsearch6.Client
	if *esToken == "" {
		es, err := elasticsearch6.NewClient(elasticsearch6.Config{
			Addresses: []string{*esAddr},
		})
		if err != nil {
			log.Fatalf("Error creating the client: %s\n", err)
		}
		esClis = append(esClis, es)
	} else {
		log.Printf("Creating OCP es client, token %s", *esToken)
		es, err := oauthEsClient(*esAddr, *esToken, *caPath, *certPath, *keyPath)
		if err != nil {
			log.Fatalf("Error creating the OCP client: %s\n", err)
		}
		esClis = append(esClis, es)

		es, err = mTlsEsClient(*esAddr, *caPath, *certPath, *keyPath)
		if err != nil {
			log.Fatalf("Error creating the OCP client: %s\n", err)
		}
		esClis = append(esClis, es)
	}

	// Get es client info
	for _, es := range esClis {
		res, err := es.Info()
		if err != nil {
			log.Fatalf("Error getting the info response: %s\n", err)
		}
		defer res.Body.Close()
		log.Printf("es URLs: %v\n", es.Transport.(*estransport.Client).URLs())

		// Get cluster version
		res, err = es.Cluster.Stats(es.Cluster.Stats.WithPretty())
		if err != nil {
			log.Fatalf("Error getting the cluster response: %s\n", err)
		}
		defer res.Body.Close()
		if res.IsError() {
			log.Printf("ERROR: %s: %s", res.Status(), res)
		} else {
			body, _ := ioutil.ReadAll(res.Body)
			str := string(body)
			version := gjson.Get(str, "nodes.versions")
			log.Println(version)
		}

		// create index
		res, err = es.Index(
			"test",                                  // Index name
			strings.NewReader(`{"title" : "Test"}`), // Document body
			es.Index.WithDocumentID("1"),            // Document ID
			es.Index.WithRefresh("true"),            // Refresh
		)
		if err != nil {
			log.Printf("ERROR: %s: %s", res.Status(), res)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("ERROR: %s: %s", res.Status(), res)
		} else {
			body, _ := ioutil.ReadAll(res.Body)
			str := string(body)
			log.Println(str)
		}
	}
}
