package handler

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	userHandlers = make(map[int]*userHandler)
}

type Handler struct {
	bot *tgbotapi.BotAPI
}

func NewHandler(bot *tgbotapi.BotAPI) *Handler {
	return &Handler{bot}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	var user *tgbotapi.User
	if update.Message != nil {
		user = update.Message.From
	} else if update.CallbackQuery != nil {
		user = update.CallbackQuery.From
	} else {
		log.Printf("Unexpected update type: %v\n", update)
		return
	}
	log.Printf("New message from user: %s\n", user.UserName)
	uh, ok := userHandlers[user.ID]
	if !ok {
		uh = &userHandler{
			user:    user,
			bot:     h.bot,
			updates: make(chan tgbotapi.Update, 100),
		}
		go uh.handleUpdates()
		userHandlers[user.ID] = uh
	}
	uh.updates <- update
}

var userHandlers map[int]*userHandler

type userHandler struct {
	user    *tgbotapi.User
	bot     *tgbotapi.BotAPI
	updates chan tgbotapi.Update
}

func (uh *userHandler) handleUpdates() {
	for update := range uh.updates {
		message := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		uh.bot.Send(message)
	}
}
