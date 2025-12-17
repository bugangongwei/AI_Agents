# AI Agents Collection

A collection of AI-powered agents built in Go, designed to perform various tasks using machine learning, APIs, and intelligent decision-making.
LangChain + 本地BGE Embedding + 本地Milvus + DeepSeek API + 国内数据源。

## Features

- HTTP REST API for agent services
- Modular agent architecture
- Integration with external APIs (weather, LLM, vector database)
- Pre-loaded knowledge base for efficient querying
- IP-based location detection (planned)

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

### Starting the Server

Build and run the HTTP server:

```bash
go build -o main.out
./main.out
```

The server will start on port 8080 and automatically load clothing rules into the vector database.

### API Endpoints

#### Outfit Recommender
Recommends daily outfits based on user questions, weather data, and clothing rules using AI.

**Endpoint:** `GET /ai_agents/outfit_recommend`

**Parameters:**
- `question` (required): User's question about outfit preferences
- `pref` (optional): Style preference (casual, formal, sporty, elegant), defaults to "casual"

**Headers:**
- `X-Forwarded-For`: Client IP address for location detection (planned feature)

**Example Request:**
```bash
curl "http://localhost:8080/ai_agents/outfit_recommend?question=What%20should%20I%20wear%20today&pref=casual"
```

**Response:**
```json
{
  "recommendation": "Based on today's weather and your preferences..."
}
```

## Project Structure

- `main.go`: HTTP server entry point and API handlers
- `outfit-recommender/`: Outfit recommender service
  - `start.go`: Core recommendation logic
  - `tools/`: Utility modules
    - `weather_tool.go`: Weather API integration
    - `llm_tool.go`: DeepSeek API client
    - `vector_tool.go`: Milvus vector database operations
    - `calendar_tool.go`: Calendar integration (legacy)
  - `data/`: Static data files
    - `clothing_rules.json`: Outfit recommendation rules
  - `.env`: Environment configuration

## Adding New Agents

1. Create a new directory for your agent.
2. Implement the agent logic in a `main.go` or similar entry file.
3. Add the agent to the `agents` map in `main.go` with name, description, usage, and path.
4. Ensure your agent handles command-line arguments appropriately.

## Dependencies

- [github.com/joho/godotenv](https://github.com/joho/godotenv) - Environment variable loading
- [github.com/sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) - OpenAI API client (used for DeepSeek)
- [github.com/milvus-io/milvus-sdk-go/v2](https://github.com/milvus-io/milvus-sdk-go) - Milvus vector database client
- QWeather API for weather data
- DeepSeek API for LLM recommendations

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.