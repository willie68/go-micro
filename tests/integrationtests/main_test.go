package integrationtests

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"testing"
)

var client http.Client

var baseURL = "https://one-dm-dev.easy02.proactcloud.de/dispatcher/gomicro"

var localExecution = flag.Bool("test.local", false, "Boolean flag to indicate whether the test should be executed against a local service.")
var skipExecution = flag.Bool("test.skip", true, "Boolean flag to indicate whether the test should be skipped.")

func TestMain(m *testing.M) {
	flag.Parse()
	log.Println(*localExecution)
	log.Println(*skipExecution)
	var exitVal = 0

	initClient()
	if *localExecution {
		baseURL = "https://127.0.0.1:8443"
	}
	exitVal = m.Run()
	os.Exit(exitVal)
}

func initClient() {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client = http.Client{
		Transport: customTransport,
	}
}


