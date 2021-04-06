package handler

import (
	"telebot/util"
	"telebot/webdriver"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

var (
	userHandlers map[int]*userHandler
	log          *zap.SugaredLogger
)

func init() {
	log = util.InitLog("handler")
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
	var text string
	var chatID int64
	var restart bool
	if update.Message != nil {
		user = update.Message.From
		text = update.Message.Text
		chatID = update.Message.Chat.ID
		if update.Message.IsCommand() {
			restart = true
		}
	} else if update.CallbackQuery != nil {
		user = update.CallbackQuery.From
		text = update.CallbackQuery.Data
		chatID = update.CallbackQuery.Message.Chat.ID
	} else {
		log.Errorw("Unexpected update type", update)
		return
	}
	log.Debugw("New message", "username", user.UserName, "message", text)
	uh, ok := userHandlers[user.ID]
	if !ok {
		uh = &userHandler{
			user:    user,
			chatID:  chatID,
			bot:     h.bot,
			updates: make(chan string, 100),
			data:    make(map[string]string),
		}
		go uh.handleUpdates()
		userHandlers[user.ID] = uh
	}
	if restart {
		uh.seqid = 0
	}
	uh.updates <- text
}

type userHandler struct {
	bot     *tgbotapi.BotAPI
	updates chan string
	user    *tgbotapi.User
	chatID  int64
	data    map[string]string
	seqid   int
}

func (uh *userHandler) handleUpdates() {
	for text := range uh.updates {
		msg := tgbotapi.NewMessage(uh.chatID, "")
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
			msg.Text = "Search may take some time, please wait"
			uh.bot.Send(msg)
			if err := webdriver.MakeAppointment(uh.data); err != nil {
				log.Errorw("hadleUpdates", "error", err)
			}
			msg.Text = "Made an appointment for: " + uh.data["timeslot"]
			msg.ReplyMarkup = results
			uh.bot.Send(msg)
		}
	}
	close(uh.updates)
	delete(userHandlers, uh.user.ID)
}
