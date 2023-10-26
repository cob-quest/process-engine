package util

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/models"
	"go.mongodb.org/mongo-driver/bson"
)

func UnmarshalJson(input []byte) map[string]interface{} {
	m := make(map[string]interface{})
	err := json.Unmarshal(input, &m)
	FailOnError(err, "Failed to convert to map")

	return m
}

func MapJsonToProcessEngine(m map[string]interface{}) *models.ProcessEngine {
	data := &models.ProcessEngine{}

	err := mapstructure.Decode(m, data)
	FailOnError(err, "Failed to convert to ProcessEngine")

	return data
}

func ConvertToBson(data map[string]interface{}) []byte {
	bsonBytes, err := bson.Marshal(data)
	FailOnError(err, "Failed to convert to BSON")

	return bsonBytes
}