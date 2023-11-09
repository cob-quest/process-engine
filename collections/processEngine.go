package collections

import (
	"context"
	"time"

	"gitlab.com/cs302-2023/g3-team8/project/process-engine/config"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/models"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/util"
	"go.mongodb.org/mongo-driver/mongo"
)

var processEngineCollection *mongo.Collection

func CreateProcessEngine(processEngine *models.ProcessEngine) (result *mongo.InsertOneResult) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if processEngineCollection == nil {
		processEngineCollection = config.OpenCollection(config.GetClient(), "process_engine")
	}

	processEngine.Timestamp = *util.GetCurrentDateTime()
	result, err := processEngineCollection.InsertOne(ctx, processEngine)

	util.FailOnError(err, "Failed to insert processEngine into db!")

	return result
}