package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsageLog struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	User      string             `json:"user" bson:"user"`
	Action    string             `json:"action" bson:"action"`
	Payload   string             `json:"payload" bson:"payload"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
}
