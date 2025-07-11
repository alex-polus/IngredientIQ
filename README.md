# IngredientIQ
🥗 AI-powered food analysis for better health decisions 🧬
## 🚀 Overview

IngredientIQ is a preventative health AI tool that helps you understand what's really in your food. By analyzing your nutrition data, we identify potentially harmful processed ingredients and predict their long-term health impacts.
### 🎯 Key Features

- **Smart Food Analysis:** Upload your nutrition data from apps like MyFitnessPal or Cronometer
- **Ingredient Detection:** AI-powered identification of ultra-processed ingredients
- **Health Impact Prediction:** Learn about potential long-term health effects
- **Personalized Insights:** Get recommendations based on your health profile
- **Real-time Chat:** Interact with our AI to ask questions about your diet

### 🧪 How It Works
1. **Upload Your Data:** Export your food log to .json from MyFitnessPal or Cronometer 
2. **AI Analysis:** Our model identifies ultra-processed ingredients
3. **Risk Assessment:** Get a detailed breakdown of potential health impacts
4. **Recommendations:** Receive personalized alternatives and suggestions

## 📦 Installation

### Prerequisites
- Go 1.22.4 or later
- OpenAI API key or compatible API (like OpenRouter)

### Setup
1. **Clone the repository:**
   ```bash
   git clone https://github.com/alexpolus/IngredientIQ.git
   cd IngredientIQ
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up environment variables (optional):**
   Create a `.env` file in the project root:
   ```
   INGRDNT_IQ_OPENAI_API_KEY=your_api_key_here
   INGRDNT_IQ_OPENAI_API_BASE_URL=your_api_base_url_here
   ```
   
   *Note: If you don't set these, the program will prompt you to enter them on first run.*

## 🚀 Usage

### Running the Application
```bash
go run main.go
```

### First Time Setup
1. **API Configuration:** If you haven't set environment variables, you'll be prompted to enter:
   - Your OpenAI API key
   - Your API base URL (e.g., for OpenRouter: `https://openrouter.ai/api/v1`)

2. **Food Log File:** Enter the path to your JSON food log file:
   - Type the full or relative path to your `.json` file
   - Example: `sample_food_log.json` or `./data/my_food_log.json`
   - The program will remember this path for future runs

### Subsequent Runs
- **Default File:** Press Enter when prompted to use your previously entered food log file
- **New File:** Type a new file path to analyze a different food log
- **Interactive Chat:** After the initial analysis, ask follow-up questions about your diet
- **Exit:** Type `quit` to exit the program

### Preparing Your Food Log
Export your food diary as JSON from:
- **MyFitnessPal:** Use third-party export tools
- **Cronometer:** Export data in JSON format
- **Manual:** Create a JSON file following the structure in `sample_food_log.json`

### Example Session
```
📁 Enter food log file path (or press Enter for default: sample_food_log.json): 
📊 Loading food log data...
✅ Food log loaded successfully!
🤖 Analyzing your food log with AI...

[AI analysis results appear here]

💬 Enter your message (or 'quit' to exit): What are the healthiest meals in my log?
🤖 Processing your question...

[AI response with specific recommendations]

💬 Enter your message (or 'quit' to exit): quit
👋 Thank you for using IngredientIQ! Stay healthy! 🌟
```