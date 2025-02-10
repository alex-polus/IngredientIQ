package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/sashabaranov/go-openai"
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
		// Verify request method and path
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check if the request headers are set correctly
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			t.Errorf("Expected Content-Type header to contain application/json, got %s", r.Header.Get("Content-Type"))
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

		// Decode and verify request body
		var reqBody openai.ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Send a mock response
		response := openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleAssistant,
						Content: "This is a test response",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create a test configuration
	config := openai.DefaultConfig("test_api_key")
	config.BaseURL = server.URL

	// Create custom transport with proper initialization
	transport := &customTransport{
		base: http.DefaultTransport,
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	client := openai.NewClientWithConfig(config)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Test message",
		},
	}

	response, err := sendRequest(client, messages)
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
