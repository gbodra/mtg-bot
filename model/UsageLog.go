package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	tb "gopkg.in/tucnak/telebot.v2"
)

type UsageLog struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Message   tb.Message         `json:"message" bson:"message"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
}
