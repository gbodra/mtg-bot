package controller

import (
	"context"
	"time"

	"github.com/gbodra/mtg-bot/model"
	"github.com/gbodra/mtg-bot/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	tb "gopkg.in/tucnak/telebot.v2"
)

var MongoClient *mongo.Client

func SaveUsageLog(m *tb.Message) {
	usageLog := model.UsageLog{
		ID:        primitive.NewObjectID(),
		Message:   *m,
		Timestamp: time.Now(),
	}

	logsCollection := MongoClient.Database("mtg").Collection("logs")

	_, err := logsCollection.InsertOne(context.TODO(), usageLog)
	utils.HandleError(err, "Error inserting usage log on Mongo")
}
