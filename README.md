# IngredientIQ
<img width="574" height="252" alt="image" src="https://github.com/user-attachments/assets/8069c677-d702-4ab0-a4cd-bb41baacf7cf" />

## 🚀 Overview

IngredientIQ is a preventative health AI tool that helps you understand what's really in your food. By analyzing your nutrition data, we identify potentially harmful processed ingredients and predict their long-term health impacts.
### 🎯 Key Features

- **Smart Food Analysis:** Upload your food nutrition data
- **Ingredient Detection:** AI-powered identification of ultra-processed ingredients
- **Health Impact Prediction:** Learn about potential long-term health effects
- **Real-time Chat:** Interact with our AI to ask questions about your diet

### 🧪 How It Works
1. **Upload Your Data:** Export your food log to .json from a food tracker (e.g. MyFitnessPal, Cronometer) or manually create one
2. **AI Analysis:** Gemini 2.5 Pro model identifies ultra-processed ingredients
3. **Risk Assessment:** Get a detailed breakdown of potential health impacts
4. **Recommendations:** Receive personalized alternatives and suggestions from Gemini

## 📦 Installation

### Prerequisites
- Go 1.20 or later
- OpenAI-compatible API endpoint (like OpenRouter, VertexAI etc.)
   - Gemini 2.5 Pro available on endpoint
  
*Note: If you don't have access to Gemini, the model can be edited [here](https://github.com/alex-polus/IngredientIQ/blob/main/main.go#L20)*

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

## 🍗 Preparing Your Food Log
To use your food diary with this project, export it as a JSON file using one of the following methods:

- **MyFitnessPal or Cronometer**:
  1.   Export your food log as a `.csv` file from the app.
  2. Simplify the `.csv` by keeping only the `food names`, `dates`, and `quantities` columns, and remove all other fields.
  3. Convert the `.csv` to JSON using an AI tool (e.g., ChatGPT, Gemini, or Grok). Aim for a structure similar to `sample_food_log.json` (exact matching is not required).
  4. *(Optional)* Use this [open-source tool](https://github.com/jrmycanady/cronometer-export) to export Cronometer logs directly to JSON.

- **Manual**:
  Create a JSON file manually, following the structure provided in `sample_food_log.json`.


## 🚀 Usage

### Running the Application
```bash
go run main.go
```

### First Time Setup
1. **API Configuration:** If you haven't set environment variables, you'll be prompted to enter:
   - Your OpenAI-compatible API key
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

### Example Session

<img width="700" height="594" alt="Screenshot 2025-07-11 at 12 46 13 PM" src="https://github.com/user-attachments/assets/d4293d7c-bbb5-44f5-8c13-065f3bc5e251" />
