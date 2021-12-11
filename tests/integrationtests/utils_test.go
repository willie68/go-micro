package integrationtests

import (
	"bytes"
	"net/http"
	"testing"
)


func executeRequest(t *testing.T, request *http.Request, status int) (http.Response, error) {
	response, err := client.Do(request)
	if err != nil {
		t.Errorf("The HTTP request failed with error %s\n", err)
		return *response, err
	} else {
		if response.StatusCode != status {
			t.Errorf("response returned wrong status code: got %v want %v", response.StatusCode, status)
		}
		expected := "application/json; charset=utf-8"
		contentType := response.Header.Get("Content-Type")
		if contentType != expected {
			t.Errorf("response returned unexpected Content-Type: got %v want %v", contentType, expected)
		}
	}
	return *response, nil
}

func createSecureRequest(t *testing.T, method string, requestURL string, requestBody []byte) *http.Request {
	request, err := http.NewRequest(method, requestURL, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Errorf("The creation of the request body failed with error %s\n", err)
	}
	request.Header.Set("X-es-apikey", "5f332fb109831d2c1e1117e3ed690d3d")
	request.Header.Set("X-es-system", "easy1")
	request.Header.Set("X-es-tenant", newTenant)
	request.Header.Set("Content-Type", "application/json")
	return request
}



