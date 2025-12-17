package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

// type EmbeddingResponse struct {
// 	Embeddings [][]float32 `json:"embeddings"`
// }

var (
	milvusClient   client.Client
	collectionName = "outfit_preferences"
)

func init() {

	// 从环境变量获取Milvus地址（Compose会自动设置）
	milvusHost := os.Getenv("MILVUS_HOST")
	if milvusHost == "" {
		milvusHost = "localhost" // 开发环境回退
	}

	milvusPort := os.Getenv("MILVUS_PORT")
	if milvusPort == "" {
		milvusPort = "19530"
	}

	address := fmt.Sprintf("%s:%s", milvusHost, milvusPort)

	log.Printf("Connecting to Milvus at %s", address)

	// Create context with timeout to prevent hanging on connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error
	milvusClient, err = client.NewClient(ctx, client.Config{
		Address: address,
	})
	if err != nil {
		log.Printf("Failed to connect to Milvus: %v", err)
	}

	// Drop collection if exists (ignore errors as collection may not exist)
	_ = milvusClient.DropCollection(ctx, collectionName)

	// Define schema for the collection
	schema := &entity.Schema{
		CollectionName: collectionName,
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     true,
			},
			{
				Name:     "text",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "1000",
				},
			},
			{
				Name:     "vector",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": "512",
				},
			},
			{
				Name:     "temperature_min",
				DataType: entity.FieldTypeInt32,
			},
			{
				Name:     "temperature_max",
				DataType: entity.FieldTypeInt32,
			},
			{
				Name:     "weather",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "50",
				},
			},
			{
				Name:     "preference",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "50",
				},
			},
			{
				Name:     "outfit",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "500",
				},
			},
		},
	}

	// Create collection
	err = milvusClient.CreateCollection(ctx, schema, 2)
	if err != nil {
		log.Printf("Failed to create collection: %v", err)
		// Note: Continue as collection might already exist or other issues
	}

	// Create index on vector field
	idx, err := entity.NewIndexFlat(entity.L2)
	if err != nil {
		log.Printf("Failed to create index: %v", err)
	} else {
		err = milvusClient.CreateIndex(ctx, collectionName, "vector", idx, false)
		if err != nil {
			log.Printf("Failed to create index on collection: %v", err)
		}
	}

	// Load collection into memory for search operations
	err = milvusClient.LoadCollection(ctx, collectionName, false)
	if err != nil {
		log.Printf("Failed to load collection: %v", err)
		// Note: Continue, but searches may fail
	}
}

func getEmbeddings(texts []string) ([][]float32, error) {
	values := url.Values{}
	for _, text := range texts {
		values.Add("texts", text)
	}
	fullURL := "http://host.docker.internal:8000/embed?" + values.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var embResp [][]float32
	err = json.Unmarshal(body, &embResp)
	if err != nil {
		return nil, err
	}

	return embResp, nil
}

func EmbedAndStore(text string) error {
	if milvusClient == nil {
		return fmt.Errorf("Milvus client not initialized")
	}

	// Get embedding
	embeddings, err := getEmbeddings([]string{text})
	if err != nil {
		return err
	}

	vec := embeddings[0]

	// Insert into Milvus
	ctx := context.Background()
	_, err = milvusClient.Insert(ctx, collectionName, "",
		entity.NewColumnVarChar("text", []string{text}),
		entity.NewColumnFloatVector("vector", 512, [][]float32{vec}),
		entity.NewColumnInt32("temperature_min", []int32{0}),
		entity.NewColumnInt32("temperature_max", []int32{0}),
		entity.NewColumnVarChar("weather", []string{""}),
		entity.NewColumnVarChar("preference", []string{""}),
		entity.NewColumnVarChar("outfit", []string{""}),
	)
	return err
}

func EmbedAndStoreRule(text string, tempMin, tempMax int, weather, preference, outfit string) error {
	if milvusClient == nil {
		return fmt.Errorf("Milvus client not initialized")
	}

	// Get embedding
	embeddings, err := getEmbeddings([]string{text})
	if err != nil {
		return err
	}
	if len(embeddings) == 0 {
		return fmt.Errorf("no embeddings for text %s returned", text)
	}

	vec := embeddings[0]

	// Insert into Milvus
	ctx := context.Background()
	_, err = milvusClient.Insert(ctx, collectionName, "",
		entity.NewColumnVarChar("text", []string{text}),
		entity.NewColumnFloatVector("vector", 512, [][]float32{vec}),
		entity.NewColumnInt32("temperature_min", []int32{int32(tempMin)}),
		entity.NewColumnInt32("temperature_max", []int32{int32(tempMax)}),
		entity.NewColumnVarChar("weather", []string{weather}),
		entity.NewColumnVarChar("preference", []string{preference}),
		entity.NewColumnVarChar("outfit", []string{outfit}),
	)
	return err
}

func StorePreference(userInput, recommendation string) error {
	text := fmt.Sprintf("User preference: %s | Recommended outfit: %s", userInput, recommendation)
	return EmbedAndStore(text)
}

func SearchSimilar(queryText string, maxTemp, minTemp float64, weather, pref string, topK int) ([]string, error) {
	if milvusClient == nil {
		return nil, fmt.Errorf("Milvus client not initialized")
	}

	// Get embedding for query
	embeddings, err := getEmbeddings([]string{queryText})
	if err != nil {
		return nil, fmt.Errorf("failed to get embeddings: %v", err)
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings for query text %s returned", queryText)
	}
	vec := embeddings[0]

	ctx := context.Background()

	// Build search request with filters
	expr := fmt.Sprintf("temperature <= %d and temperature >= %d and weather == '%s' and preference == '%s'", int(maxTemp), int(minTemp), weather, pref)

	sp, err := entity.NewIndexFlatSearchParam()
	if err != nil {
		return nil, fmt.Errorf("failed to create search param: %v", err)
	}

	sr, err := milvusClient.Search(ctx, collectionName, []string{}, expr, []string{"outfit"}, []entity.Vector{entity.FloatVector(vec)}, "vector", entity.L2, topK, sp)
	if err != nil {
		return nil, fmt.Errorf("search failed: %v", err)
	}

	var results []string
	for _, result := range sr {
		for _, fieldData := range result.Fields {
			if outfitCol, ok := fieldData.(*entity.ColumnVarChar); ok {
				results = append(results, outfitCol.Data()...)
			}
		}
	}

	return results, nil
}
