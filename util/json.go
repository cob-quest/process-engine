package util

import (
	"encoding/json"
	"log"

	"gitlab.com/cs302-2023/g3-team8/project/process-engine/models"
	"go.mongodb.org/mongo-driver/bson"
)

func UnmarshalJson(input []byte) *models.ProcessEngine {
	data := &models.ProcessEngine{}
	log.Println("Input string:", input)
	err := json.Unmarshal(input, &data)
	FailOnError(err, "Failed to convert to map")
	return data
}

func ConvertToBson(data map[string]interface{}) []byte {
	bsonBytes, err := bson.Marshal(data)
	FailOnError(err, "Failed to convert to BSON")

	return bsonBytes
}