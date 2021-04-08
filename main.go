package main

import (
	"log"
	"telebot/handler"
	"telebot/util"
	"telebot/webdriver"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	logger := util.InitLog("main")
	config := util.InitConfig()
	webdriver.InitDriver(config)

	bot, err := tgbotapi.NewBotAPI(config.GetString("apitoken"))
	if err != nil {
		logger.Panic(err)
	}
	logger.Debugw("Authorized", "accountname", bot.Self.UserName)

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
