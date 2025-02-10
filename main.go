package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

const (
	customBaseURL = "https://openrouter.ai/api/v1"
	modelName     = "deepseek/deepseek-r1-distill-llama-70b"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("API key not found in .env file")
	}

	// Create custom configuration
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = customBaseURL

	// Create custom transport with proper initialization
	transport := &customTransport{
		base: http.DefaultTransport,
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	// Initialize the client with custom configuration
	client := openai.NewClientWithConfig(config)

	foodLog, err := readFoodLog("sample_food_log.json")
	if err != nil {
		log.Fatalf("Error reading food log: %v", err)
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a prevanative health expert analyzing an included daily food log. Identify potentially harmful processed foods and predict long-term health impacts. List the unhealthiest processed foods in the log and provide insights and recommendations to improve long term health outcomes based on findings and analysis.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Analyze this food log and provide insights: %s", foodLog),
		},
	}

	response, err := sendRequest(client, messages)
	if err != nil {
		log.Fatalf("Error sending request to API: %v", err)
	}

	fmt.Println("AI Response:")
	fmt.Println(response)

	// Start interactive chat
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nEnter your message (or 'quit' to exit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "quit" {
			break
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		})

		response, err := sendRequest(client, messages)
		if err != nil {
			log.Printf("Error sending request to API: %v", err)
			continue
		}

		fmt.Println("AI Response:")
		fmt.Println(response)
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: response,
		})
	}
}

func readFoodLog(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func sendRequest(client *openai.Client, messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    modelName,
			Messages: messages,
		},
	)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return resp.Choices[0].Message.Content, nil
}

// customTransport adds required headers for OpenRouter
type customTransport struct {
	base http.RoundTripper
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.base == nil {
		t.base = http.DefaultTransport
	}

	// Ensure Content-Type is set
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add OpenRouter specific headers
	req.Header.Set("HTTP-Referer", "https://github.com/alexpolus/IngredientIQ")
	req.Header.Set("X-Title", "IngredientIQ")

	return t.base.RoundTrip(req)
}
