package controller

import (
	"context"

	"github.com/gbodra/mtg-bot/model"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

func SaveUsageLog(usageLog model.UsageLog) {
	logsCollection := MongoClient.Database("mtg").Collection("logs")

	_, err := logsCollection.InsertOne(context.TODO(), usageLog)
	if err != nil {
		panic(err)
	}
}
