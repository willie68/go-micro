package integrationtests

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
)

var newTrait string
var newType string
var newDisplay string

var newTenant = "testtenant"

func TestCarpenterEndpoints(t *testing.T) {
	if *skipExecution {
		t.Skip("Integration tests are skipped by default")
	}
	t.Run("Test1: checkGetReadiness", checkGetReadiness)
	t.Run("Test2: checkGetHealthCheck", checkGetHealthCheck)
}



func checkGetReadiness(t *testing.T) {
	type Ready struct {
		Message string `json:"message"`
		Name    string `json:"name"`
	}
	requestURL := baseURL + "/health/readiness"
	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		t.Errorf("The creation of the request  failed with error %s\n", err)
	}

	response, err := executeRequest(t, request, http.StatusOK)
	if err == nil {
		expected := "service started"
		decoder := json.NewDecoder(response.Body)
		val := &Ready{}
		err := decoder.Decode(val)
		if err != nil {
			log.Fatal(err)
		}
		if val.Message != expected {
			t.Errorf("response returned unexpected message: got %v want %v",
				val.Message, expected)
		}
	}
}

func checkGetHealthCheck(t *testing.T) {
	type Health struct {
		Message   string `json:"message"`
		LastCheck string `json:"lastCheck"`
	}

	requestURL := baseURL + "/health/health"
	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		t.Errorf("The creation of the request  failed with error %s\n", err)
	}
	response, err := executeRequest(t, request, http.StatusOK)
	if err == nil {
		expected := "service up and running"
		decoder := json.NewDecoder(response.Body)
		val := &Health{}
		err := decoder.Decode(val)
		if err != nil {
			log.Fatal(err)
		}
		if val.Message != expected {
			t.Errorf("response returned unexpected message: got %v want %v",
				val.Message, expected)
			log.Println(val.LastCheck)
		}
	}
}
