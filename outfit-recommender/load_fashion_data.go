package outfit_recommender

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"outfit-recommender/outfit-recommender/tools"
)

func LoadFashionTrend() {
	file, err := os.Open("outfit-recommender/data/Fashion Dataset.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		text := fmt.Sprintf("Year: %s, Collection: %s, Type: %s, Location: %s, Gender: %s, Colour: %s, Patterns: %s, Shoes: %s, Accessories: %s, Fabric: %s, Theme: %s",
			record[0], record[1], record[2], record[3], record[4], record[5], record[6], record[7], record[8], record[9], record[10])
		err := tools.EmbedAndStore(text)
		if err != nil {
			log.Printf("Error storing entry %d: %v", i, err)
		}
	}

	fmt.Println("Data loaded successfully")
}
