package models

import (
	"google.golang.org/genproto/googleapis/type/datetime"
)

type ProcessEngine struct {
	Timestamp    datetime.DateTime `json:"timestamp" bson:"timestamp"`
	CorId        *string           `json:"corId" bson:"corId"`
	Event        *string           `json:"event" bson:"event"`
	EventSuccess bool              `json:"eventSuccess" bson:"eventSuccess"`
	TestCreator  *string           `json:"testCreator" bson:"testCreator"`
	ImageName    *string           `json:"imageName" bson:"imageName"`
	Participants []string          `json:"participants" bson:"participants"`
}
