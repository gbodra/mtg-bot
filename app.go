package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gbodra/mtg-bot/controller"
	"github.com/gbodra/mtg-bot/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	tb "gopkg.in/tucnak/telebot.v2"
)

type App struct {
	Port  string
	Bot   *tb.Bot
	Mongo *mongo.Client
}

func (a *App) Initialize() {
	err := godotenv.Load()
	utils.HandleError(err, "Error loading .env file")

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	a.Mongo, err = mongo.Connect(context.TODO(), clientOptions)
	utils.HandleError(err, "Error connecting to Mongo")

	a.injectClients()
	a.initializeTelegramBot()
}

func (a *App) injectClients() {
	controller.MongoClient = a.Mongo
}

func (a *App) Run() {
	a.telegramBotActions()
}

func (a *App) initializeTelegramBot() {
	log.Println("Telegram bot initiated...")

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TG_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	utils.HandleError(err, "Error initializing Telegram Bot")

	a.Bot = b
}
