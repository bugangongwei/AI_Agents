# AI Agents Collection

A collection of AI-powered agents built in Go, designed to perform various tasks using machine learning, APIs, and intelligent decision-making.

## Project Structure

- `tools/`: Shared utility tools that can be used by different agents.
  - `calendar_tool.go`: Google Calendar API integration.
  - `llm_tool.go`: LLM (DeepSeek) integration.
  - `location_tool.go`: IP-based city resolution.
  - `vector_tool.go`: Milvus vector database operations.
  - `weather_tool.go`: Weather API integration.
- `outfit_recommender/`: An agent that provides clothing recommendations based on weather, location, and schedule.
  - `data/`: Contains data files like `IP2LOCATION-LITE-DB5_CN.CSV` and `clothing_rules.json`.
- `.env`: Environment variables for API keys and configuration.

## Features

- **Shared Tools**: Reusable components for common tasks (AI, Weather, etc.).
- **Schedule Integration**: Automatically fetches user schedules from Google Calendar.
- **Smart Location**: Resolves city names from user IP using a local database.
- **Vector Search**: Uses Milvus for high-performance retrieval of clothing rules.
- **LLM Powered**: Generates natural language recommendations using DeepSeek.

## Getting Started

### Prerequisites

- Go 1.24+
- Docker & Docker Compose (for Milvus)
- API Keys:
  - `DEEPSEEK_API_KEY`: For LLM recommendations.
  - `WEATHER_API_TOKEN`: For weather data.
  - `GOOGLE_CALENDAR_CREDENTIALS`: Path to your Google Calendar credentials JSON.

### Configuration

Fill in your API keys in the `.env` file at the root.

### Running with Docker Compose

This is the recommended way to start all services including Milvus:

```bash
docker compose up -d
```

### Running the API Server

```bash
go run main.go
```

The server will be available at `http://localhost:8081`.

### API Endpoint

**GET `/ai_agents/outfit_recommend`**

Parameters:
- `question`: (Required) User's question (e.g., "What should I wear today?").
- `pref`: (Optional) Style preference (e.g., "casual", "business").
- `loc`: (Optional) City name. If omitted, resolved from IP.
- `schedule`: (Optional) Manual schedule description. If omitted, fetched from Google Calendar.

Example:
```
http://localhost:8081/ai_agents/outfit_recommend?question=明天穿什么&pref=casual
```
