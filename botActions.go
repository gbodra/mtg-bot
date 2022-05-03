package main

import (
	"time"

	"github.com/gbodra/mtg-bot/controller"
	"github.com/go-co-op/gocron"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (a *App) telegramBotActions() {
	a.Bot.Handle(tb.OnText, func(m *tb.Message) {
		controller.SaveUsageLog(m)
		a.Bot.Send(m.Sender, "Sorry your request was invalid. Please check the menu to see what are the available actions")
	})

	a.Bot.Handle("/hello", func(m *tb.Message) {
		controller.HelloWorld(m, a.Bot)
	})

	a.Bot.Handle("/find_card_by_id", func(m *tb.Message) {
		controller.GetCardInfoById(m, a.Bot)
	})

	a.Bot.Handle("/find_card_by_name", func(m *tb.Message) {
		controller.GetCardInfoByName(m, a.Bot)
	})

	a.Bot.Handle("/alert_optin", func(m *tb.Message) {
		controller.AlertOptin(m, a.Bot)
	})

	a.Bot.Handle("/alert_optout", func(m *tb.Message) {
		controller.AlertOptout(m, a.Bot)
	})

	scheduler := gocron.NewScheduler(time.Local)
	scheduler.Every("30s").Do(func() {
		// scheduler.Every(1).Day().At("10:22").Do(func() {
		controller.SendNotification(a.Bot)
		// recipient, _ := a.Bot.ChatByID("605145454")
		// a.Bot.Send(recipient, "Test cron")
	})
	scheduler.StartAsync()

	a.Bot.Start()
}
