package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBalancer(t *testing.T) {
	// TODO: Реалізуйте юніт-тест для балансувальникка.
	serversPool = []Server{
		{name: "server1:8080", counter: 10},
		{name: "server2:8080", counter: 5},
		{name: "server3:8080", counter: 7},
	}

	// Act
	leastLoadedServer := findLeastLoadedServer()

	// Assert
	if leastLoadedServer.name != "server2:8080" {
		t.Errorf("Expected server2:8080, but got %s", leastLoadedServer.name)
	}
	if leastLoadedServer.counter != 5 {
		t.Errorf("Expected counter 5, but got %d", leastLoadedServer.counter)
	}
}

func TestForward(t *testing.T) {
	// Arrange
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer mockServer.Close()

	server := Server{name: mockServer.Listener.Addr().String(), counter: 0}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Act
	err := forward(server.name, w, req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, but got %v", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "OK" {
		t.Errorf("Expected body to be 'OK', but got %s", body)
	}
}

func TestHealth(t *testing.T) {
	// Arrange
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	// Act
	isHealthy := health(mockServer.Listener.Addr().String())

	// Assert
	if !isHealthy {
		t.Errorf("Expected server to be healthy, but it is not")
	}
}
