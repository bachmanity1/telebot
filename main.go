package main

import (
	"telebot/handler"
	"telebot/scraper"
	"telebot/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	log := util.InitLog("main")
	defer util.Recover(log)
	config := util.InitConfig()

	// init booth data
	if err := scraper.InitData(config); err != nil {
		log.Panicw("InitData", "error", err)
	}

	// init bot
	bot, err := tgbotapi.NewBotAPI(config.GetString("apitoken"))
	if err != nil {
		log.Panicw("InitBot", "error", err)
	}
	log.Debugw("Authorized", "accountname", bot.Self.UserName)
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
