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
	var messageID int
	if update.Message != nil {
		user = update.Message.From
		text = update.Message.Text
		chatID = update.Message.Chat.ID
		messageID = update.Message.MessageID
		if update.Message.Command() == "start" {
			uh, ok := userHandlers[user.ID]
			if ok {
				uh.done <- true
				<-uh.events
			}
			uh = &userHandler{
				user:        user,
				chatID:      chatID,
				bot:         h.bot,
				events:      make(chan *event, 100),
				requestData: make(map[string]string),
				expectedID:  messageID,
				done:        make(chan bool),
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
		h.bot.Send(tgbotapi.NewMessage(chatID, "Type /start to make an appointment"))
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
	done          chan bool
}

type event struct {
	value string
	id    int
}

func (uh *userHandler) handleEvents() {
	defer util.Recover(log)
	wdDone := make(chan bool)
	go func() {
		<-uh.done
		close(wdDone)
		delete(userHandlers, uh.user.ID)
		close(uh.events)
		log.Debugw("handleEvents cleanup done", "username", uh.user.UserName)
	}()
	for event := range uh.events {
		if uh.expectedID == event.id {
			if event.value == "exit" {
				uh.bot.Send(tgbotapi.NewMessage(uh.chatID, "Bye Bye!"))
				break
			}
			uh.requestData[uh.expectedField] = event.value
			nextMessage := uh.getNextMessage()
			if uh.expectedField == "receipt" {
				uh.bot.Send(tgbotapi.NewMessage(uh.chatID, "Search may take some time(potentially hours!), we'll send you a message as soon as we have an update for you"))
				receipt, err := webdriver.MakeAppointment(uh.requestData, wdDone)
				if err != nil {
					log.Errorw("Make Appointment", "error", err)
					uh.bot.Send(tgbotapi.NewMessage(uh.chatID, "Something went wrong (probably unavailable booth), please retry"))
					uh.replyID = 0
					nextMessage = uh.getNextMessage()
				} else {
					nextMessage.Text = receipt
					go webdriver.CancelPrevAppointment(uh.requestData)
				}
			}
			if uh.expectedField == "booth" {
				boothes, err := webdriver.GetBoothes(uh.requestData)
				if err != nil {
					log.Errorw("GetBoothes", "error", err)
				}
				if err != nil {
					uh.bot.Send(tgbotapi.NewMessage(uh.chatID, "Wrong username or password, please retry"))
					uh.replyID = 0
					nextMessage = uh.getNextMessage()
				} else {
					nextMessage.ReplyMarkup = makeBoothMarkup(boothes)
				}
			}
			msg, _ := uh.bot.Send(nextMessage)
			log.Debugw("Send Message", "text", msg.Text, "messageID", msg.MessageID)
			uh.expectedID = msg.MessageID
			if nextMessage.ReplyMarkup == nil {
				uh.expectedID++
			}
		}
	}
	log.Debugw("handleUpdate exit", "username", uh.user.UserName)
}

func (uh *userHandler) getNextMessage() tgbotapi.MessageConfig {
	reply := replies[len(replies)-1]
	if uh.replyID < len(replies) {
		reply = replies[uh.replyID]
	}
	message := tgbotapi.NewMessage(uh.chatID, reply.text)
	if reply.isMarkup {
		message.ReplyMarkup = reply.markup
	}
	uh.expectedField = reply.field
	uh.replyID++
	return message
}
