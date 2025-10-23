package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
)

const (
	modelName  = "google/gemini-2.5-pro"
	configFile = ".ingredientiq_config"
)

// renderMarkdown renders markdown text with proper terminal formatting
func renderMarkdown(text string) string {
	// Render markdown with terminal width of 80 characters and indent of 2
	return string(markdown.Render(text, 80, 2))
}

// readSecureInput reads input from terminal without echoing characters (for passwords/API keys)
// Falls back to regular input if not in an interactive terminal
func readSecureInput() (string, error) {
	// Check if stdin is a terminal
	if term.IsTerminal(int(syscall.Stdin)) {
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		fmt.Println() // Add newline after password input
		return strings.TrimSpace(string(bytePassword)), nil
	} else {
		// Fallback to regular input for non-interactive terminals (e.g., piped input)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(input), nil
	}
}

// loadLastFilePath loads the last used file path from config file
func loadLastFilePath() string {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return "" // Return empty string if config file doesn't exist
	}
	return strings.TrimSpace(string(content))
}

// saveLastFilePath saves the file path to config file
func saveLastFilePath(filePath string) error {
	return ioutil.WriteFile(configFile, []byte(filePath), 0644)
}

// promptForFilePath prompts user for food log file path with default option
func promptForFilePath(reader *bufio.Reader) (string, error) {
	lastPath := loadLastFilePath()

	if lastPath != "" {
		color.New(color.FgCyan, color.Bold).Printf("📁 Enter food log file path (or press Enter for default: %s): ", lastPath)
	} else {
		color.New(color.FgCyan, color.Bold).Print("📁 Enter food log file path: ")
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// If user pressed Enter without typing anything, use last path
	if input == "" && lastPath != "" {
		input = lastPath
		color.New(color.FgGreen, color.Bold).Printf("✅ Using default file: %s\n", lastPath)
	} else if input == "" {
		return "", fmt.Errorf("no file path provided")
	}

	// Check if file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", input)
	}

	// Save the file path for next time
	err := saveLastFilePath(input)
	if err != nil {
		color.New(color.FgYellow, color.Bold).Printf("⚠️  Warning: Could not save file path for next time: %v\n", err)
	}

	return input, nil
}

// displayLogo shows a colorful ASCII logo for IngredientIQ
func displayLogo() {
	// Create color functions
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)
	red := color.New(color.FgRed, color.Bold)

	fmt.Println()
	green.Println("  ██╗███╗   ██╗ ██████╗ ██████╗ ███████╗██████╗ ██╗███████╗███╗   ██╗████████╗")
	green.Println("  ██║████╗  ██║██╔════╝ ██╔══██╗██╔════╝██╔══██╗██║██╔════╝████╗  ██║╚══██╔══╝")
	green.Println("  ██║██╔██╗ ██║██║  ███╗██████╔╝█████╗  ██║  ██║██║█████╗  ██╔██╗ ██║   ██║   ")
	green.Println("  ██║██║╚██╗██║██║   ██║██╔══██╗██╔══╝  ██║  ██║██║██╔══╝  ██║╚██╗██║   ██║   ")
	green.Println("  ██║██║ ╚████║╚██████╔╝██║  ██║███████╗██████╔╝██║███████╗██║ ╚████║   ██║   ")
	green.Println("  ╚═╝╚═╝  ╚═══╝ ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═════╝ ╚═╝╚══════╝╚═╝  ╚═══╝   ╚═╝   ")
	fmt.Println()
	magenta.Print("                              ██╗ ██████╗ ")
	red.Println("")
	magenta.Print("                              ██║██╔═══██╗")
	red.Println("")
	magenta.Print("                              ██║██║   ██║")
	red.Println("")
	magenta.Print("                              ██║██║▄▄ ██║")
	red.Println("")
	magenta.Print("                              ██║╚██████╔╝")
	red.Println("")
	magenta.Print("                              ╚═╝ ╚══▀▀═╝ ")
	red.Println("")
	fmt.Println()
	cyan.Println("              🥗 AI-Powered Nutrition Analysis for Better Health 🧬")
	fmt.Println()
	color.New(color.FgWhite, color.Bold).Println("═══════════════════════════════════════════════════════════════════════════════")
	fmt.Println()
}

func main() {
	// Display colorful logo
	displayLogo()

	// Start interactive chat
	reader := bufio.NewReader(os.Stdin)
	godotenv.Overload(".env")
	apiKey := os.Getenv("INGRDNT_IQ_OPENAI_API_KEY")
	baseURL := os.Getenv("INGRDNT_IQ_OPENAI_API_BASE_URL")
	err := error(nil)
	if apiKey == "" || baseURL == "" {
		apiKey, baseURL, err = setApiKeyAndURL(*reader)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create custom configuration
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL

	// Create custom transport with proper initialization
	transport := &customTransport{
		base: http.DefaultTransport,
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	// Initialize the client with custom configuration
	client := openai.NewClientWithConfig(config)

	// Prompt for food log file path
	filePath, err := promptForFilePath(reader)
	if err != nil {
		color.New(color.FgRed, color.Bold).Printf("❌ Error getting file path: %v\n", err)
		log.Fatalf("Error getting file path: %v", err)
	}
	fmt.Println()

	// Read food log with colored status message
	color.New(color.FgCyan, color.Bold).Println("📊 Loading food log data...")
	foodLog, err := readFoodLog(filePath)
	if err != nil {
		color.New(color.FgRed, color.Bold).Printf("❌ Error reading food log: %v\n", err)
		log.Fatalf("Error reading food log: %v", err)
	}

	systemPrompt, err := readFoodLog("system_prompt.md")
	if err != nil {
		log.Fatalf("Error reading system prompt: %v", err)
	}
	color.New(color.FgGreen, color.Bold).Println("✅ Food log loaded successfully!")
	fmt.Println()

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf("%s", systemPrompt),
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Analyze this food log and provide insights: %s", foodLog),
		},
	}

	// Send request to AI with colored status
	color.New(color.FgYellow, color.Bold).Println("🤖 Analyzing your food log with AI...")
	response, err := sendRequest(client, messages)
	if err != nil {
		color.New(color.FgRed, color.Bold).Printf("❌ Error sending request to API: %v\n", err)
		log.Errorf("Error sending request to API: %v", err)
		color.New(color.FgYellow).Print("\n🔧 Clear API key / URL? (y/n): ")
		edit, _ := reader.ReadString('\n')
		edit = strings.TrimSpace(edit)
		if edit == "y" {
			err = godotenv.Write(map[string]string{"INGRDNT_IQ_OPENAI_API_KEY": "", "INGRDNT_IQ_OPENAI_API_BASE_URL": ""}, ".env")
			if err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(1)
	}

	fmt.Println()
	color.New(color.FgGreen, color.Bold).Println("🎯 AI HEALTH ANALYSIS RESULTS:")
	color.New(color.FgWhite, color.Bold).Println("═══════════════════════════════════════════════════════════════════════════════")
	fmt.Println()
	// Render markdown response with proper formatting
	formattedResponse := renderMarkdown(response)
	fmt.Print(formattedResponse)
	fmt.Println()
	color.New(color.FgWhite, color.Bold).Println("═══════════════════════════════════════════════════════════════════════════════")

	for {
		fmt.Println()
		color.New(color.FgCyan, color.Bold).Print("💬 Enter your message (or 'quit' to exit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "quit" {
			color.New(color.FgMagenta, color.Bold).Println("\n👋 Thank you for using IngredientIQ! Stay healthy! 🌟")
			break
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		})

		color.New(color.FgYellow, color.Bold).Println("\n🤖 Processing your question...")
		response, err := sendRequest(client, messages)
		if err != nil {
			color.New(color.FgRed, color.Bold).Printf("❌ Error sending request to API: %v\n", err)
			log.Printf("Error sending request to API: %v", err)
			continue
		}

		fmt.Println()
		color.New(color.FgGreen, color.Bold).Println("🤖 AI RESPONSE:")
		color.New(color.FgWhite, color.Bold).Println("───────────────────────────────────────────────────────────────────────────────")
		// Render markdown response with proper formatting
		formattedResponse := renderMarkdown(response)
		fmt.Print(formattedResponse)
		color.New(color.FgWhite, color.Bold).Println("───────────────────────────────────────────────────────────────────────────────")
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: response,
		})
	}
}

func setApiKeyAndURL(reader bufio.Reader) (apiKey string, baseURL string, err error) {
	apiKey = os.Getenv("INGRDNT_IQ_OPENAI_API_KEY")
	if apiKey == "" {
		color.New(color.FgYellow, color.Bold).Print("🔑 Enter your API key (input hidden): ")
		apiKey, err = readSecureInput()
		if err != nil {
			color.New(color.FgRed, color.Bold).Printf("❌ Error reading API key: %v\n", err)
			return "", "", err
		}

		err := os.Setenv("INGRDNT_IQ_OPENAI_API_KEY", apiKey)
		if err != nil {
			color.New(color.FgRed, color.Bold).Println("❌ API key failed to save to .env")
			log.Fatal("API key failed to save to .env")
		}
		color.New(color.FgGreen, color.Bold).Println("✅ API key saved successfully!")
	}

	baseURL = os.Getenv("INGRDNT_IQ_OPENAI_API_BASE_URL")
	if baseURL == "" {
		color.New(color.FgYellow, color.Bold).Print("🌐 Enter your OpenAI API base URL: ")
		baseURL, _ = reader.ReadString('\n')
		baseURL = strings.TrimSpace(baseURL)

		err := os.Setenv("INGRDNT_IQ_OPENAI_BASE_URL", baseURL)
		if err != nil {
			color.New(color.FgRed, color.Bold).Println("❌ Base URL failed to save to .env")
			log.Fatal("API key failed to save to .env")
		}
		color.New(color.FgGreen, color.Bold).Println("✅ Base URL saved successfully!")
	}
	err = godotenv.Write(map[string]string{"INGRDNT_IQ_OPENAI_API_KEY": apiKey, "INGRDNT_IQ_OPENAI_API_BASE_URL": baseURL}, ".env")
	if err != nil {
		log.Fatal(err)
	}
	return apiKey, baseURL, nil
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
