package main

import (
	"telebot/handler"
	"telebot/util"
	"telebot/webdriver"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	log := util.InitLog("main")
	config := util.InitConfig()

	bot, err := tgbotapi.NewBotAPI(config.GetString("apitoken"))
	if err != nil {
		log.Panic(err)
	}
	log.Debugw("Authorized", "accountname", bot.Self.UserName)

	defer util.Recover(log)
	webdriver.InitDriver(config)
	h := handler.NewHandler(bot)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		h.HandleUpdate(update)
	}

}
