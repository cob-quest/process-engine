package models

import (
	"google.golang.org/genproto/googleapis/type/datetime"
)

type ProcessEngine struct {
	Timestamp    datetime.DateTime `json:"timestamp" bson:"timestamp"`
	CorId        *string           `json:"corId" bson:"corId"`
	Event        *string           `json:"event" bson:"event"`
	EventSuccess bool              `json:"eventSuccess" bson:"eventSuccess"`
	CreatorName  *string           `json:"creatorName" bson:"creatorName,omitempty"`
	ImageName    *string           `json:"imageName" bson:"imageName,omitempty"`
	Participants *[]string         `json:"participants" bson:"participants,omitempty"`
}
