package handler

import (
	"telebot/scraper"
	"telebot/util"

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
	var uType updateType
	if update.Message != nil {
		uType = plainText
		user = update.Message.From
		text = update.Message.Text
		chatID = update.Message.Chat.ID
		messageID = update.Message.MessageID
		if update.Message.Command() == "start" {
			uh, ok := userHandlers[user.ID]
			if ok {
				uh.done <- true
				<-uh.done
			}
			uh = &userHandler{
				user:        user,
				chatID:      chatID,
				bot:         h.bot,
				events:      make(chan *event, 100),
				requestData: make(map[string]string),
				nextID:      messageID,
				done:        make(chan bool),
				isActive:    true,
			}
			go uh.handleEvents()
			userHandlers[user.ID] = uh
		}
	} else if update.CallbackQuery != nil {
		uType = callbackQuery
		user = update.CallbackQuery.From
		text = update.CallbackQuery.Data
		chatID = update.CallbackQuery.Message.Chat.ID
		messageID = update.CallbackQuery.Message.MessageID
	} else {
		log.Errorw("Unnext update type", "update", update)
		return
	}
	log.Debugw("Received message", "username", user.UserName, "text", text, "messageID", messageID)
	uh, ok := userHandlers[user.ID]
	if !ok {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Type /start to make an appointment"))
		return
	}
	uh.events <- &event{text, messageID, uType}
}

type userHandler struct {
	bot         *tgbotapi.BotAPI
	user        *tgbotapi.User
	chatID      int64
	events      chan *event
	requestData map[string]string
	nextID      int
	nextType    updateType
	nextField   string
	replyID     int
	done        chan bool
	isActive    bool
}

type event struct {
	value string
	id    int
	uType updateType
}

func (uh *userHandler) handleEvents() {
	defer util.Recover(log)
	wdDone := make(chan bool)
	go func() {
		<-uh.done
		uh.isActive = false
		delete(userHandlers, uh.user.ID)
		close(uh.events)
		close(uh.done)
		close(wdDone)
		log.Debugw("handleEvents cleanup", "username", uh.user.UserName)
	}()
	for event := range uh.events {
		if uh.isValidEvent(event) {
			if event.value == "exit" {
				uh.send(tgbotapi.NewMessage(uh.chatID, "Bye Bye!"))
				uh.done <- true
				break
			}
			uh.requestData[uh.nextField] = event.value
			nextMessage := uh.getNextMessage()
			if uh.nextField == "receipt" {
				uh.send(tgbotapi.NewMessage(uh.chatID, "Search may take some time(potentially hours!), we'll send you a message as soon as we have an update for you"))
				receipt, err := scraper.MakeAppointment(uh.requestData)
				if err != nil {
					log.Errorw("Make Appointment", "error", err)
					uh.send(tgbotapi.NewMessage(uh.chatID, "Something went wrong (probably unavailable booth), please retry"))
					uh.replyID = 0
					nextMessage = uh.getNextMessage()
				} else {
					nextMessage.Text = receipt
				}
			}
			if uh.nextField == "booth" {
				nextMessage.ReplyMarkup = boothMarkup[uh.requestData["branch"]]
			}
			msg, _ := uh.send(nextMessage)
			log.Debugw("Send Message", "text", msg.Text, "messageID", msg.MessageID)
			uh.nextType = plainText
			if nextMessage.ReplyMarkup != nil {
				uh.nextID = msg.MessageID
				uh.nextType = callbackQuery
			}
		}
	}
	log.Debugw("handleEvents exit", "username", uh.user.UserName)
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
	uh.nextField = reply.field
	uh.replyID++
	return message
}

func (uh *userHandler) send(message tgbotapi.MessageConfig) (tgbotapi.Message, error) {
	if uh.isActive {
		return uh.bot.Send(message)
	}
	return tgbotapi.Message{}, nil
}

func (uh *userHandler) isValidEvent(e *event) bool {
	if e.uType == plainText {
		return e.uType == uh.nextType
	}
	return e.uType == uh.nextType && e.id == uh.nextID
}
