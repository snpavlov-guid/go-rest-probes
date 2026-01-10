// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"net/http"
	"os"
	"testing"

	"github.com/imroc/req/v3"
)

var (
	baseURL = flag.String("base-url", getEnvOrDefault("AIRCRAFT_API_BASE_URL", "http://localhost:8080"),
		"Base URL for the aircraft API service")
	client = req.C()
)

// getEnvOrDefault retrieves an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestBasicEndpoints tests basic endpoints that don't require authentication
func TestGreet(t *testing.T) {
	url := *baseURL + "/"

	resp, err := client.R().Get(url)
	if err != nil {
		t.Fatalf("Failed to make GET request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestGreetWithText(t *testing.T) {
	testText := "world"
	url := *baseURL + "/" + testText

	resp, err := client.R().Get(url)
	if err != nil {
		t.Fatalf("Failed to make GET request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestVersion(t *testing.T) {
	url := *baseURL + "/version"

	resp, err := client.R().Get(url)
	if err != nil {
		t.Fatalf("Failed to make GET request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestProfile_Unauthenticated(t *testing.T) {
	url := *baseURL + "/profile"

	resp, err := client.R().Get(url)
	if err != nil {
		t.Fatalf("Failed to make GET request to %s: %v", url, err)
	}

	// Profile endpoint may return 200 with a message indicating user is not authenticated
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

// TestAircraftEndpoints tests aircraft-related endpoints
func TestGetAircrafts(t *testing.T) {
	url := *baseURL + "/api/v1/aircrafts"

	// Test without query parameters
	resp, err := client.R().Get(url)
	if err != nil {
		t.Fatalf("Failed to make GET request to %s: %v", url, err)
	}

	// Note: This endpoint requires authentication, so we expect 401 or potentially 200 if JWT is handled
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

func TestGetAircrafts_WithPagination(t *testing.T) {
	url := *baseURL + "/api/v1/aircrafts"

	resp, err := client.R().
		SetQueryParam("size", "10").
		SetQueryParam("offset", "0").
		Get(url)
	if err != nil {
		t.Fatalf("Failed to make GET request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

func TestGetAircraftByCode(t *testing.T) {
	testCode := "SU9"
	url := *baseURL + "/api/v1/aircrafts/" + testCode

	resp, err := client.R().Get(url)
	if err != nil {
		t.Fatalf("Failed to make GET request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

func TestGetAircraftByCode_EmptyCode(t *testing.T) {
	url := *baseURL + "/api/v1/aircrafts/"

	resp, err := client.R().Get(url)
	if err != nil {
		t.Fatalf("Failed to make GET request to %s: %v", url, err)
	}

	// Empty code should result in 404 or 500
	if resp.StatusCode < http.StatusBadRequest && resp.StatusCode != http.StatusNotFound {
		t.Logf("Note: Empty code returned status %d (expected 400, 404, or 500)", resp.StatusCode)
	}
}

func TestCreateAircraft(t *testing.T) {
	url := *baseURL + "/api/v1/aircrafts/create"

	aircraftInput := map[string]interface{}{
		"code":   "TEST001",
		"nameRu": "Тестовый самолет",
		"nameEn": "Test Aircraft",
		"range":  5000,
	}

	resp, err := client.R().
		SetBodyJsonMarshal(aircraftInput).
		Post(url)
	if err != nil {
		t.Fatalf("Failed to make POST request to %s: %v", url, err)
	}

	// Expected 401 (unauthorized) since this requires auth, or 200 if authenticated
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 200, 400, or 401, got %d", resp.StatusCode)
	}
}

func TestUpdateAircraft(t *testing.T) {
	url := *baseURL + "/api/v1/aircrafts/update"

	aircraftInput := map[string]interface{}{
		"code":   "TEST001",
		"nameRu": "Обновленный самолет",
		"nameEn": "Updated Aircraft",
		"range":  6000,
	}

	resp, err := client.R().
		SetBodyJsonMarshal(aircraftInput).
		Post(url)
	if err != nil {
		t.Fatalf("Failed to make POST request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 200, 400, or 401, got %d", resp.StatusCode)
	}
}

func TestDeleteAircraft_POST(t *testing.T) {
	testCode := "TEST001"
	url := *baseURL + "/api/v1/aircrafts/delete/" + testCode

	resp, err := client.R().Post(url)
	if err != nil {
		t.Fatalf("Failed to make POST request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

func TestDeleteAircraft_DELETE(t *testing.T) {
	testCode := "TEST001"
	url := *baseURL + "/api/v1/aircrafts/" + testCode

	resp, err := client.R().Delete(url)
	if err != nil {
		t.Fatalf("Failed to make DELETE request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

// TestAirportEndpoints tests airport-related endpoints
func TestGetAirports(t *testing.T) {
	url := *baseURL + "/api/v1/airports"

	resp, err := client.R().Get(url)
	if err != nil {
		t.Fatalf("Failed to make GET request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

func TestGetAirports_WithPagination(t *testing.T) {
	url := *baseURL + "/api/v1/airports"

	resp, err := client.R().
		SetQueryParam("size", "10").
		SetQueryParam("offset", "0").
		Get(url)
	if err != nil {
		t.Fatalf("Failed to make GET request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

func TestGetAirportByCode(t *testing.T) {
	testCode := "SVO"
	url := *baseURL + "/api/v1/airports/" + testCode

	resp, err := client.R().Get(url)
	if err != nil {
		t.Fatalf("Failed to make GET request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

func TestCreateAirport(t *testing.T) {
	url := *baseURL + "/api/v1/airports/create"

	airportInput := map[string]interface{}{
		"code":     "TEST",
		"nameRu":   "Тестовый аэропорт",
		"nameEn":   "Test Airport",
		"cityRu":   "Тестовый город",
		"cityEn":   "Test City",
		"timezone": "UTC+3",
	}

	resp, err := client.R().
		SetBodyJsonMarshal(airportInput).
		Post(url)
	if err != nil {
		t.Fatalf("Failed to make POST request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 200, 400, or 401, got %d", resp.StatusCode)
	}
}

func TestUpdateAirport(t *testing.T) {
	url := *baseURL + "/api/v1/airports/update"

	airportInput := map[string]interface{}{
		"code":     "TEST",
		"nameRu":   "Обновленный аэропорт",
		"nameEn":   "Updated Airport",
		"cityRu":   "Обновленный город",
		"cityEn":   "Updated City",
		"timezone": "UTC+4",
	}

	resp, err := client.R().
		SetBodyJsonMarshal(airportInput).
		Post(url)
	if err != nil {
		t.Fatalf("Failed to make POST request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 200, 400, or 401, got %d", resp.StatusCode)
	}
}

func TestDeleteAirport_POST(t *testing.T) {
	testCode := "TEST"
	url := *baseURL + "/api/v1/airports/delete/" + testCode

	resp, err := client.R().Post(url)
	if err != nil {
		t.Fatalf("Failed to make POST request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

func TestDeleteAirport_DELETE(t *testing.T) {
	testCode := "TEST"
	url := *baseURL + "/api/v1/airports/" + testCode

	resp, err := client.R().Delete(url)
	if err != nil {
		t.Fatalf("Failed to make DELETE request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

// Helper function to create authenticated requests
func createAuthenticatedRequest(token string, body interface{}) *req.Request {
	r := client.R()

	if token != "" {
		r.SetHeader("Authorization", "Bearer "+token)
	}

	if body != nil {
		r.SetBodyJsonMarshal(body)
	}

	return r
}

// TestAuthenticatedEndpoint demonstrates how to test with authentication
func TestAuthenticatedEndpoint_Example(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping authenticated endpoint test in short mode")
	}

	url := *baseURL + "/api/v1/aircrafts"

	// For actual authenticated tests, you would need a valid JWT token
	// token := "your-jwt-token-here"
	token := os.Getenv("AIRCRAFT_API_TOKEN") // Get token from environment

	resp, err := createAuthenticatedRequest(token, nil).Get(url)
	if err != nil {
		t.Fatalf("Failed to make authenticated request: %v", err)
	}

	// With valid token, expect 200
	// Without token, expect 401
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}
