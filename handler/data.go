package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

const replyLen = 7

var seqidToData = map[int]string{
	1: "branch",
	2: "nationality",
	3: "purpose",
	4: "username",
	5: "password",
	6: "phone",
}

var branches = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Seoul", "1"),
		tgbotapi.NewInlineKeyboardButtonData("Daejon", "2"),
		tgbotapi.NewInlineKeyboardButtonData("Busan", "3"),
	),
)

var nationalities = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Chinese", "1"),
		tgbotapi.NewInlineKeyboardButtonData("non-Chinese", "2"),
	),
)

var purposes = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Application", "1"),
		tgbotapi.NewInlineKeyboardButtonData("Extension", "2"),
		tgbotapi.NewInlineKeyboardButtonData("Change", "3"),
	),
)

var results = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Good Enough!", "exit"),
		tgbotapi.NewInlineKeyboardButtonData("Look for earlier date", "2"),
	),
)

type reply struct {
	text     string
	isMarkup bool
	markup   tgbotapi.InlineKeyboardMarkup
}

var seqidToReplies = map[int]*reply{
	0: {text: "Choose Immigration Branch", isMarkup: true, markup: branches},
	1: {text: "Choose your nationality", isMarkup: true, markup: nationalities},
	2: {text: "Choose purpose of visit", isMarkup: true, markup: purposes},
	3: {text: "Enter your username", isMarkup: false},
	4: {text: "Enter your password", isMarkup: false},
	5: {text: "Enter your phone number", isMarkup: false},
}
