package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gbodra/mtg-bot/model"
	"github.com/gbodra/mtg-bot/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

func HelloWorld(m *tb.Message, bot *tb.Bot) {
	bot.Send(m.Sender, "Hello I'm Miyagi your friend who will keep you up to date with your MTG card prices")

	SaveUsageLog(m)
}

func GetCardInfoById(m *tb.Message, bot *tb.Bot) {
	id := m.Payload

	response, err := http.Get(os.Getenv("API_URI") + "/cards/" + id)
	utils.HandleError(err, "Error getting the alert from the API on GetCardInfoById")

	sendCardInfo(response.Body, m, bot)
}

func GetCardInfoByName(m *tb.Message, bot *tb.Bot) {
	name := m.Payload

	response, err := http.Get(os.Getenv("API_URI") + "/cards?q=" + url.QueryEscape(name))
	utils.HandleError(err, "Error getting the alert from the API on GetCardInfoByName")

	sendCardInfo(response.Body, m, bot)
}

func sendCardInfo(body io.ReadCloser, m *tb.Message, bot *tb.Bot) {
	responseData, err := ioutil.ReadAll(body)
	utils.HandleError(err, "Error reading the response on sendCardInfo")

	var cardObject model.Card
	json.Unmarshal(responseData, &cardObject)

	colorsJson, err := json.Marshal(cardObject.Colors)
	utils.HandleError(err, "Error transforming object into json sendCardInfo")

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
		message += "   Date: " + el.CreatedAt.Format(time.RFC822) + "\n"
		message += "-----------------------------\n"
	}

	message += cardObject.ScryfallURL

	bot.Send(m.Sender, message)

	SaveUsageLog(m)
}
