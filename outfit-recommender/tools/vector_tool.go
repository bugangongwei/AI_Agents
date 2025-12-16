package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type EmbeddingResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

var (
	milvusClient   client.Client
	collectionName = "outfit_preferences"
)

func init() {
	ctx := context.Background()
	c, err := client.NewClient(ctx, client.Config{
		Address: "localhost:19530",
	})
	if err != nil {
		log.Printf("Failed to connect to Milvus: %v", err)
		return
	}
	milvusClient = c

	// Create collection if not exists
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
					"dim": "768",
				},
			},
		},
	}
	err = milvusClient.CreateCollection(ctx, schema, 2)
	if err != nil {
		log.Printf("Failed to create collection: %v", err)
	}
}

func getEmbeddings(texts []string) ([][]float32, error) {
	reqBody := map[string][]string{"texts": texts}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:8000/embed", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var embResp EmbeddingResponse
	err = json.Unmarshal(body, &embResp)
	if err != nil {
		return nil, err
	}

	return embResp.Embeddings, nil
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
		entity.NewColumnFloatVector("vector", 768, [][]float32{vec}),
	)
	return err
}

func StorePreference(userInput, recommendation string) error {
	text := fmt.Sprintf("User preference: %s | Recommended outfit: %s", userInput, recommendation)
	return EmbedAndStore(text)
}

func SearchSimilar(userInput string, topK int) ([]string, error) {
	// Temporarily return empty list for testing
	return []string{}, nil
}
