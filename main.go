package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	defaultAPIURL = "https://openrouter.ai/api/v1/chat/completions"
	model         = "openai/gpt-3.5-turbo"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type APIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type APIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

var apiURL = defaultAPIURL

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("API key not found in .env file")
	}

	foodLog, err := readFoodLog("food_log.json")
	if err != nil {
		log.Fatalf("Error reading food log: %v", err)
	}

	messages := []Message{
		{Role: "system", Content: "You are a nutritionist analyzing a food log. Provide insights and recommendations based on the log."},
		{Role: "user", Content: fmt.Sprintf("Analyze this food log and provide insights: %s", foodLog)},
	}

	response, err := sendRequest(apiKey, messages)
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

		messages = append(messages, Message{Role: "user", Content: input})
		response, err := sendRequest(apiKey, messages)
		if err != nil {
			log.Printf("Error sending request to API: %v", err)
			continue
		}

		fmt.Println("AI Response:")
		fmt.Println(response)
		messages = append(messages, Message{Role: "assistant", Content: response})
	}
}

func readFoodLog(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func sendRequest(apiKey string, messages []Message) (string, error) {
	requestBody, err := json.Marshal(APIRequest{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(requestBody)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/alexpolus/IngredientIQ")
	req.Header.Set("X-Title", "IngredientIQ")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", err
	}

	if len(apiResponse.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return apiResponse.Choices[0].Message.Content, nil
}
