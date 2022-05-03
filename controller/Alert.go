package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gbodra/mtg-bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
)

// TODO: criar logica para buscar precos que tiveram grande mudanca

func AlertOptin(m *tb.Message, bot *tb.Bot) {
	_, err := http.Post(os.Getenv("API_URI")+"/alert?chat_id="+strconv.FormatInt(m.Chat.ID, 10), "", nil)

	if err != nil {
		log.Println(err)
		bot.Send(m.Sender, "Error saving your alert")
	} else {
		bot.Send(m.Sender, "Alert saved successfully!")
	}

	SaveUsageLog(m)
}

func AlertOptout(m *tb.Message, bot *tb.Bot) {
	chatId := strconv.FormatInt(m.Chat.ID, 10)
	alertId := getAlert(chatId)
	client := &http.Client{}

	url := os.Getenv("API_URI") + "/alert/" + alertId
	request, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		log.Println(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		bot.Send(m.Sender, "Error deleting your alert")
	} else {
		bot.Send(m.Sender, "You no longer will receive price alerts")
	}

	defer resp.Body.Close()

	SaveUsageLog(m)
}

func getAlert(chatId string) string {
	response, err := http.Get(os.Getenv("API_URI") + "/alert/" + chatId)

	if err != nil {
		log.Println(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	var alertObject model.Alert
	json.Unmarshal(responseData, &alertObject)

	return alertObject.ID.Hex()
}

func getAlerts() []model.Alert {
	response, err := http.Get(os.Getenv("API_URI") + "/listAlerts")
	if err != nil {
		log.Println(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	var alertObject []model.Alert
	json.Unmarshal(responseData, &alertObject)

	return alertObject
}

func getAlertMessage() model.AlertMessage {
	response, err := http.Get(os.Getenv("API_URI") + "/price/top")
	if err != nil {
		log.Println(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	var alertMessageObject model.AlertMessage
	json.Unmarshal(responseData, &alertMessageObject)

	return alertMessageObject
}

func SendNotification(bot *tb.Bot) {
	listOfRecipients := getAlerts()
	message := getAlertMessage()
	messageStr := ""

	for _, card := range message.Cards {
		messageStr += "Card Name: " + card.Name + "\n"
		messageStr += "Current Price: US$" + fmt.Sprintf("%.2f", card.LastPrice) + "\n"
		messageStr += "Change from last day: US$" + fmt.Sprintf("%.2f", card.NormalMovementMoney) + "\n"
		messageStr += "Change from last day: " + fmt.Sprintf("%.2f", card.NormalMovementPercentage*100) + "%\n"
		messageStr += "-----------------------------\n\n"
	}
	log.Println("Alerta deveria ser enviado")

	for _, el := range listOfRecipients {
		recipient, _ := bot.ChatByID(strconv.FormatInt(el.ChatID, 10))
		bot.Send(recipient, messageStr)
	}
}
