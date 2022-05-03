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
	"github.com/gbodra/mtg-bot/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

func AlertOptin(m *tb.Message, bot *tb.Bot) {
	_, err := http.Post(os.Getenv("API_URI")+"/alert?chat_id="+strconv.FormatInt(m.Chat.ID, 10), "", nil)

	if err != nil {
		utils.HandleError(err, "Error on saving the optin")
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
	utils.HandleError(err, "Error creating the DELETE request on AlertOptout")

	resp, err := client.Do(request)
	if err != nil {
		utils.HandleError(err, "Error executing the DELETE request on AlertOptout")
		bot.Send(m.Sender, "Error deleting your alert")
	} else {
		bot.Send(m.Sender, "You no longer will receive price alerts")
	}

	defer resp.Body.Close()

	SaveUsageLog(m)
}

func getAlert(chatId string) string {
	response, err := http.Get(os.Getenv("API_URI") + "/alert/" + chatId)
	utils.HandleError(err, "Error getting the alert from the API on getAlert")

	responseData, err := ioutil.ReadAll(response.Body)
	utils.HandleError(err, "Error reading the response on getAlert")

	var alertObject model.Alert
	json.Unmarshal(responseData, &alertObject)

	return alertObject.ID.Hex()
}

func getAlerts() []model.Alert {
	response, err := http.Get(os.Getenv("API_URI") + "/listAlerts")
	utils.HandleError(err, "Error getting the alert from the API on getAlerts")

	responseData, err := ioutil.ReadAll(response.Body)
	utils.HandleError(err, "Error reading the response on getAlerts")

	var alertObject []model.Alert
	json.Unmarshal(responseData, &alertObject)

	return alertObject
}

func getAlertMessage() model.AlertMessage {
	response, err := http.Get(os.Getenv("API_URI") + "/price/top")
	utils.HandleError(err, "Error getting the alert from the API on getAlertMessage")

	responseData, err := ioutil.ReadAll(response.Body)
	utils.HandleError(err, "Error reading the response on getAlertMessage")

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
		recipient, err := bot.ChatByID(strconv.FormatInt(el.ChatID, 10))
		utils.HandleError(err, "Error getting the chatId")
		bot.Send(recipient, messageStr)
	}
}
