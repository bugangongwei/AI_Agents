package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	outfit_recommender "outfit-recommender/outfit-recommender"
)

// TODO: Add IP to location parsing with ip2region

func outfitRecommendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	question := r.URL.Query().Get("question")
	pref := r.URL.Query().Get("pref")

	if question == "" {
		http.Error(w, "Missing 'question' parameter", http.StatusBadRequest)
		return
	}
	if pref == "" {
		pref = "casual"
	}

	// TODO: Parse location from IP
	location := "Shanghai" // Default

	args := []string{"-question", question, "-pref", pref, "-loc", location}

	recommendation, err := outfit_recommender.Start(args)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"recommendation": recommendation}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Load clothing rules into Milvus on startup
	err := outfit_recommender.LoadClothingRules()
	if err != nil {
		log.Printf("Failed to load clothing rules: %v", err)
	}

	http.HandleFunc("/ai_agents/outfit_recommend", outfitRecommendHandler)

	port := "8080"
	fmt.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
