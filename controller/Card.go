package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gbodra/mtg-bot/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	tb "gopkg.in/tucnak/telebot.v2"
)

func HelloWorld(m *tb.Message, bot *tb.Bot) {
	bot.Send(m.Sender, "Hello I'm Miyagi your friend who will keep you up to date with your MTG card prices")

	usageLog := model.UsageLog{
		ID:        primitive.NewObjectID(),
		User:      m.Sender.Username,
		Action:    "/hello",
		Timestamp: time.Now(),
		Payload:   m.Payload,
	}
	SaveUsageLog(usageLog)
}

func GetCardInfo(m *tb.Message, bot *tb.Bot) {
	id := m.Payload

	response, err := http.Get(os.Getenv("API_URI") + "/cards/" + id)

	if err != nil {
		log.Println(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	var cardObject model.Card
	json.Unmarshal(responseData, &cardObject)

	colorsJson, _ := json.Marshal(cardObject.Colors)

	message := "Name: " + cardObject.Name + "\n"
	message += "Type: " + cardObject.TypeLine + "\n"
	message += "Colors: " + string(colorsJson) + "\n"
	message += "Rarity: " + cardObject.Rarity + "\n"
	message += "Set: " + cardObject.Set + "\n"
	message += "Set Name: " + cardObject.SetName + "\n"
	message += "Price:\n"

	for i, el := range cardObject.Prices.Prices {
		message += "   [" + strconv.Itoa(i) + "]\n"
		message += "   Type: " + el.PrintingType + "\n"
		message += "   Market Price: " + fmt.Sprintf("%.2f", el.MarketPrice) + "\n"
		message += "   Buy List Market Price: " + fmt.Sprintf("%.2f", el.BuylistMarketPrice) + "\n"
		message += "   Listed Median Price: " + fmt.Sprintf("%.2f", el.ListedMedianPrice) + "\n"
		message += "-----------------------------\n"
	}

	message += cardObject.ScryfallURL

	bot.Send(m.Sender, message)

	usageLog := model.UsageLog{
		ID:        primitive.NewObjectID(),
		User:      m.Sender.Username,
		Action:    "/find_card_by_id",
		Timestamp: time.Now(),
		Payload:   m.Payload,
	}
	SaveUsageLog(usageLog)
}
