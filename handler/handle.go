package handler

import (
	"telebot/util"
	"telebot/webdriver"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

const startcmd = "start"

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
	var messageID int
	if update.Message != nil {
		user = update.Message.From
		text = update.Message.Text
		chatID = update.Message.Chat.ID
		messageID = update.Message.MessageID
		if update.Message.Command() == startcmd {
			uh, ok := userHandlers[user.ID]
			if ok {
				delete(userHandlers, user.ID)
				close(uh.events)
			}
			uh = &userHandler{
				user:        user,
				chatID:      chatID,
				bot:         h.bot,
				events:      make(chan *event, 100),
				requestData: make(map[string]string),
				expectedID:  messageID,
			}
			go uh.handleEvents()
			userHandlers[user.ID] = uh
		}
	} else if update.CallbackQuery != nil {
		user = update.CallbackQuery.From
		text = update.CallbackQuery.Data
		chatID = update.CallbackQuery.Message.Chat.ID
		messageID = update.CallbackQuery.Message.MessageID
	} else {
		log.Errorw("Unexpected update type", "update", update)
		return
	}
	log.Debugw("Received message", "username", user.UserName, "text", text, "messageID", messageID)
	uh, ok := userHandlers[user.ID]
	if !ok {
		h.bot.Send(tgbotapi.NewMessage(chatID, "type /start to make an appointment"))
		return
	}
	uh.events <- &event{text, messageID}
}

type userHandler struct {
	bot           *tgbotapi.BotAPI
	user          *tgbotapi.User
	chatID        int64
	events        chan *event
	requestData   map[string]string
	expectedID    int
	expectedField string
	replyID       int
}

type event struct {
	value string
	id    int
}

func (uh *userHandler) handleEvents() {
	for event := range uh.events {
		if uh.expectedID == event.id {
			if event.value == "exit" {
				uh.bot.Send(tgbotapi.NewMessage(uh.chatID, "Bye Bye!"))
				break
			}
			uh.requestData[uh.expectedField] = event.value
			if uh.expectedField == "branch" {
				subbranches, err := webdriver.GetSubBranches(uh.requestData)
				if err != nil {
					uh.bot.Send(tgbotapi.NewMessage(uh.chatID, "Wrong username or password, please retry"))
					uh.replyID = 0
				}
				insertSubbranchMarkup(subbranches)
			}
			nextMessage, ok := uh.getNextMessage()
			if !ok {
				timeslot, err := webdriver.MakeAppointment(uh.requestData)
				uh.requestData["prevtimeslot"] = timeslot
				if err != nil {
					log.Errorw("Make Appointment", "error", err)
				}
				nextMessage = tgbotapi.NewMessage(uh.chatID, "Made an appointment for: "+timeslot)
				nextMessage.ReplyMarkup = results
			}
			msg, _ := uh.bot.Send(nextMessage)
			log.Debugw("Send Message", "text", msg.Text, "messageID", msg.MessageID)
			uh.expectedID = msg.MessageID
			if nextMessage.ReplyMarkup == nil {
				uh.expectedID++
			}
		}
	}
}

func (uh *userHandler) getNextMessage() (tgbotapi.MessageConfig, bool) {
	if uh.replyID >= len(replies) {
		uh.expectedField = ""
		return tgbotapi.MessageConfig{}, false
	}
	reply := replies[uh.replyID]
	message := tgbotapi.NewMessage(uh.chatID, reply.text)
	if reply.isMarkup {
		message.ReplyMarkup = reply.markup
	}
	uh.expectedField = reply.field
	uh.replyID++
	return message, true
}
