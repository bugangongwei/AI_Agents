# AI Agents Collection

A collection of AI-powered agents built in Go, designed to perform various tasks using machine learning, APIs, and intelligent decision-making. Includes an outfit recommender optimized for Chinese language support.
LangChain + 本地BGE Embedding + 本地Milvus + DeepSeek API + 国内数据源。

## Features

- HTTP REST API for agent services
- Modular agent architecture
- Integration with external APIs (weather, LLM, vector database)
- Pre-loaded knowledge base for efficient querying
- IP-based location detection (planned)

## Usage

### Starting the Server

#### Option 1: Run with Go

Build and run the HTTP server:

```bash
go build -o main.out
./main.out
```

The server will start on port 8081 and automatically load clothing rules into the vector database.

#### Option 2: Run with Docker Compose

Start all services using Docker Compose (includes Milvus, etcd, MinIO, and the AI Agents API):

```bash
#!/bin/bash
# Improved Docker Compose startup script with error handling

set -e  # Exit on any error

# Configuration
COMPOSE_FILE="docker-compose.yml"
PROJECT_NAME="ai-agents"

# Function to check if Docker and Docker Compose are available
check_dependencies() {
    if ! command -v docker &> /dev/null; then
        echo "Error: Docker is not installed. Please install Docker first."
        exit 1
    fi

    if ! docker info &> /dev/null; then
        echo "Error: Docker daemon is not running. Please start Docker."
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        echo "Error: Docker Compose is not installed. Please install Docker Compose."
        exit 1
    fi
}

# Function to start services
start_services() {
    echo "Starting AI Agents services with Docker Compose..."

    # Use docker compose (new syntax) if available, fallback to docker-compose
    if docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
    else
        COMPOSE_CMD="docker-compose"
    fi

    # Stop any existing services first
    $COMPOSE_CMD -f "$COMPOSE_FILE" down || true

    # Build and start services
    if ! $COMPOSE_CMD -f "$COMPOSE_FILE" up -d --build; then
        echo "Error: Failed to start services with Docker Compose."
        exit 1
    fi

    echo "Services started successfully!"
    echo "AI Agents API: http://localhost:8081"
    echo "Milvus: localhost:19530"
    echo "MinIO Console: http://localhost:9091 (admin/admin)"
}

# Function to show status
show_status() {
    echo "Service Status:"
    if docker compose version &> /dev/null; then
        docker compose -f "$COMPOSE_FILE" ps
    else
        docker-compose -f "$COMPOSE_FILE" ps
    fi
}

# Main execution
check_dependencies
start_services
show_status

echo "All services are running. Use 'docker-compose logs -f' to monitor logs."
```

**Prerequisites:**
- Docker and Docker Compose installed
- Embedding API running on port 8000 (if needed separately)

**Notes:**
- This starts the complete environment including Milvus vector database
- Data persists in the `./volumes` directory
- Use `docker-compose down` to stop all services

### API Endpoints

#### Outfit Recommender
Recommends daily outfits based on user questions, weather data, and clothing rules using AI. Optimized for Chinese language inputs and outputs.

**Endpoint:** `GET /ai_agents/outfit_recommend`

**Parameters:**
- `question` (required): User's question about outfit preferences (in Chinese)
- `pref` (optional): Style preference (休闲, 正式, 运动, 优雅), defaults to "休闲"

**Headers:**
- `X-Forwarded-For`: Client IP address for location detection (planned feature)

**Example Request:**
```bash
curl "http://localhost:8081/ai_agents/outfit_recommend?question=今天我应该穿什么&pref=休闲"
```

**Response:**
```json
{"recommendation":"今天天气晴朗，气温在15到25°C之间，早晚温差较大，建议你采用“多层穿搭法”，方便根据温度变化调整衣物。\n\n**推荐搭配：**\n- **上身**：内搭一件长袖T恤或薄款针织衫，外穿一件轻薄夹克或风衣。中午升温时可以脱掉外套，早晚凉爽时穿上保暖。\n- **下身**：选择长裤（如牛仔裤、休闲裤）即可，透气材质更适合晴天活动。\n- **鞋袜**：轻便的运动鞋或休闲鞋，搭配普通棉袜。\n\n**决策理由：**\n1. **天气条件参考**：晴天日照充足，中午体感较暖，但早晚温度偏低（15°C左右），多层穿搭能灵活应对温差。\n2. **用户需求**：问题简洁，侧重日常实用性，因此推荐通用且舒适的搭配。\n3. **穿衣规则**：遵循了“温差大时分层穿戴”的原则，兼顾保暖与透气性。\n\n**小建议**：如果白天需要长时间在户外，可以戴一顶遮阳帽或准备薄围巾，防晒的同时也能防风。祝你今天过得舒适！"}
```

#### Embedding API
Provides text embeddings using local BGE model.

**Endpoint:** `GET /embed`

**Parameters:**
- `texts` (required, multiple): Text strings to embed

**Example Request:**
```bash
curl "http://localhost:8000/embed?texts=hello&texts=world"
```

**Response:**
```json
{
  "embeddings": [[0.1, 0.2, ...], [0.3, 0.4, ...]]
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

````
# 1. 进入项目目录
cd /path/to/your/project

# 2. 启动所有服务（后台运行）
docker-compose up -d

# 3. 查看运行状态
docker-compose ps

# 4. 查看日志
docker-compose logs           # 所有服务日志
docker-compose logs myapp     # 只看myapp日志
docker-compose logs -f milvus # 实时查看milvus日志
# 查看完整的错误日志
docker logs --tail 100 milvus-standalone
# 或者使用docker-compose
docker-compose logs --tail 100 milvus

# 5. 停止服务
docker-compose stop           # 停止但不删除容器
docker-compose down           # 停止并删除容器
docker-compose down -v        # 停止并删除容器和卷

# 6. 重启服务
docker-compose restart

# 7. 重建并启动（代码更新后）
docker-compose up -d --build

# 8. 进入容器
docker-compose exec myapp sh  # 进入myapp容器
docker-compose exec milvus bash  # 进入milvus容器
````