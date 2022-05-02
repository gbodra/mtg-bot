package controller

import (
	"encoding/json"
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
	alertId := getAlerts(chatId)
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

func getAlerts(chatId string) string {
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
