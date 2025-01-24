package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestReadFoodLog(t *testing.T) {
	// Create a temporary food log file for testing
	tempFile, err := os.CreateTemp("", "test_food_log.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	testLog := `[{"date": "2025-01-23", "meals": [{"type": "Breakfast", "items": [{"name": "Oatmeal", "quantity": "1 cup"}]}]}]`
	if _, err := tempFile.Write([]byte(testLog)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Test reading the food log
	content, err := readFoodLog(tempFile.Name())
	if err != nil {
		t.Fatalf("readFoodLog failed: %v", err)
	}

	if content != testLog {
		t.Errorf("readFoodLog returned unexpected content. Got %s, want %s", content, testLog)
	}
}

func TestSendRequest(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request headers are set correctly
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be application/json")
		}
		if r.Header.Get("Authorization") != "Bearer test_api_key" {
			t.Errorf("Expected Authorization header to be Bearer test_api_key")
		}
		if r.Header.Get("HTTP-Referer") != "https://github.com/alexpolus/IngredientIQ" {
			t.Errorf("Expected HTTP-Referer header to be https://github.com/alexpolus/IngredientIQ")
		}
		if r.Header.Get("X-Title") != "IngredientIQ" {
			t.Errorf("Expected X-Title header to be IngredientIQ")
		}

		// Send a mock response
		response := APIResponse{
			Choices: []struct {
				Message Message `json:"message"`
			}{
				{
					Message: Message{
						Role:    "assistant",
						Content: "This is a test response",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Override the apiURL for testing
	originalAPIURL := apiURL
	apiURL = server.URL
	defer func() { apiURL = originalAPIURL }()

	messages := []Message{
		{Role: "user", Content: "Test message"},
	}

	response, err := sendRequest("test_api_key", messages)
	if err != nil {
		t.Fatalf("sendRequest failed: %v", err)
	}

	expectedResponse := "This is a test response"
	if response != expectedResponse {
		t.Errorf("sendRequest returned unexpected response. Got %s, want %s", response, expectedResponse)
	}
}

func TestMain(m *testing.M) {
	// Set up test environment
	os.Setenv("OPENAI_API_KEY", "test_api_key")

	// Run tests
	exitCode := m.Run()

	// Clean up
	os.Unsetenv("OPENAI_API_KEY")

	os.Exit(exitCode)
}
