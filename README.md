# AI Agents Collection

A collection of AI-powered agents built in Go, designed to perform various tasks using machine learning, APIs, and intelligent decision-making.
LangChain + 本地BGE Embedding + 本地Milvus + DeepSeek API + 国内数据源。

## Features

- Modular agent architecture
- Easy-to-extend framework for adding new agents
- Command-line interface for running agents
- Integration with external APIs (weather, calendar, LLM, etc.)

## Installation

1. Ensure you have Go 1.24 or later installed.
2. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/AI_Agents.git
   cd AI_Agents
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```

## Usage

Run the main program with an agent name and its arguments:

```bash
go run . <agent-name> [args...]
```

For help:
```bash
go run . --help
```

### Available Agents

#### Outfit Recommender
Recommends daily outfits based on weather, schedule, user preferences, and input using AI.

**Usage:**
```bash
go run . outfit-recommender -input "<user description>" -pref <style> -loc "<city>"
```

- `-input`: Text describing your preferences (e.g., "I like blue colors")
- `-pref`: Style preference (casual, formal, sporty, elegant)
- `-loc`: City name for weather data (e.g., "Beijing")

## Project Structure

- `main.go`: Entry point and agent dispatcher
- `outfit-recommender/`: Outfit recommender agent
  - `start.go`: Main logic
  - `tools/`: Utility tools (weather, calendar, LLM, vector)
  - `data/`: Data files (clothing rules)

## Adding New Agents

1. Create a new directory for your agent.
2. Implement the agent logic in a `main.go` or similar entry file.
3. Add the agent to the `agents` map in `main.go` with name, description, usage, and path.
4. Ensure your agent handles command-line arguments appropriately.

## Dependencies

- [github.com/joho/godotenv](https://github.com/joho/godotenv) - Environment variable loading
- [github.com/sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) - OpenAI API client
- [github.com/tmc/langchaingo](https://github.com/tmc/langchaingo) - LangChain for Go
- Google APIs for calendar and other services

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.