package handler

import (
	"log"
	"time"

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
			data:    make(map[string]string),
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
	data    map[string]string
	seqid   int
}

func (uh *userHandler) handleUpdates() {
	for update := range uh.updates {
		var text string
		var chatID int64
		if update.Message != nil {
			text = update.Message.Text
			chatID = update.Message.Chat.ID
			if update.Message.IsCommand() {
				uh.seqid = 0
			}
		} else {
			text = update.CallbackQuery.Data
			chatID = update.CallbackQuery.Message.Chat.ID
		}
		msg := tgbotapi.NewMessage(chatID, "")
		if text == "exit" {
			msg.Text = "Bye Bye!"
			uh.bot.Send(msg)
			break
		}
		uh.data[seqidToData[uh.seqid]] = text
		replydata := seqidToReplies[uh.seqid]
		if replydata != nil {
			msg.Text = replydata.text
			if replydata.isMarkup {
				msg.ReplyMarkup = replydata.markup
			}
			uh.bot.Send(msg)
			uh.seqid++
		} else {
			msg.Text = "Search may take some, plase wait"
			uh.bot.Send(msg)
			time.Sleep(3 * time.Second)
			msg.Text = "2021/03/31"
			msg.ReplyMarkup = results
			uh.bot.Send(msg)
		}
	}
	close(uh.updates)
	delete(userHandlers, uh.user.ID)
	log.Println("----------------------------------------------------------")
}
