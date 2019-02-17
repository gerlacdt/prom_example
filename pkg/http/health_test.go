package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("ERROR creating http-recorder, %v", err)
	}

	h := &HealthHandler{}
	h.ServeHTTP(w, req)
	resp := w.Result()

	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Fatalf("ERROR closing body, %v", err)
		}
	}()
	if err != nil {
		t.Fatalf("ERROR read http response, %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode, got: %d, expected: %d", resp.StatusCode, 200)
	}
	var healthResponse HealthResponse
	err = json.Unmarshal(body, &healthResponse)
	if err != nil {
		t.Fatalf("ERROR json unmarshalling: %v", err)
	}
	if healthResponse.Status != "OK" {
		t.Errorf("health status, got %s, expected: %s", healthResponse.Status, "OK")
	}
}
