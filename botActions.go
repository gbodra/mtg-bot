package main

import (
	"github.com/gbodra/mtg-bot/controller"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (a *App) telegramBotActions() {
	a.Bot.Handle("/hello", func(m *tb.Message) {
		controller.HelloWorld(m, a.Bot)
	})

	a.Bot.Handle("/find_card_by_id", func(m *tb.Message) {
		controller.GetCardInfo(m, a.Bot)
	})

	a.Bot.Start()
}