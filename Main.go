// This file contains an example how to use Infinispan HTTP/2 interface in the simplest possible way.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"crypto/tls"
	"golang.org/x/net/http2"
	"flag"
	"crypto/x509"
	"strings"
)

var (
	infinispanAddress = flag.String("address", "https://localhost:8443", "Infinispan address")
	certificatePath = flag.String("certificate", "./certificate.pem", "Certificate for accessing Infinispan")
)

func main() {
	flag.Parse()

	caCert, err := ioutil.ReadFile(*certificatePath)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// We need to create a custom transport with proper certificate and HTTP/2
	tlsTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			// If you can't want to deal with certificates just uncomment the line below.
			//InsecureSkipVerify: true,
			RootCAs:      caCertPool,
		},
	}

	// This call is essential. In case of custom transport make sure you will call it.
	http2.ConfigureTransport(tlsTransport)

	// We are creating a new HTTP Client with configured HTTP/2 transport.
	http2Client := http.Client{
		Transport: tlsTransport,
	}

	// Let's put something in cache
	url := fmt.Sprintf("%v/rest/default/test", *infinispanAddress)
	response, err := http2Client.Post(url, "text/plain", strings.NewReader("Infinispan is cool!"))
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != 200 {
		log.Printf("Status is different than 200 (was %v)!", response.StatusCode)
	}

	// And now let's get it back
	response, err = http2Client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: (%v), Protocol %v", string(body), response.Proto)
}
