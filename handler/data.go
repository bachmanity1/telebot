package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

var seqidToData = map[int]string{
	1: "branch",
	2: "purpose",
	3: "username",
	4: "password",
	5: "phone",
}

var branches = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Seoul Immigration Office", "1270667"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Sejongno Branch Office", "1271020"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Incheon Immigration Office", "1270700"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Incheon Immigration Office Ansan Branch Office", "1272143"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Suwon Immigration Office", "1270947"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Ulsan  Immigration Office", "1270698"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Busan Immigration Office", "1270686"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Daejeon Immigration Office", "1270727"),
	),
)

var purposes = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Foreign Resident Registration", "F01"),
	), tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Reissue of Alien Registration Card", "F02"),
	), tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Visa extension", "F03"),
	), tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Change of visa status", "F04"),
	), tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Granting a visa", "F05"),
	), tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Permit for other activities beyond current visa status", "F06"),
	), tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Change/addition of workplace", "F07"),
	), tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Re-entry permit (single/multiple)", "F08"),
	), tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Change of residence", "F09"),
	), tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Change of Registration Matters (Passport Information)", "F10"),
	),
)

var results = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Good Enough!", "exit"),
	), tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Look for an earlier date", "2"),
	),
)

type reply struct {
	text     string
	isMarkup bool
	markup   tgbotapi.InlineKeyboardMarkup
}

var seqidToReplies = map[int]*reply{
	0: {text: "Choose Immigration Branch", isMarkup: true, markup: branches},
	1: {text: "Choose purpose of visit", isMarkup: true, markup: purposes},
	2: {text: "Enter your username", isMarkup: false},
	3: {text: "Enter your password", isMarkup: false},
	4: {text: "Enter your phone number", isMarkup: false},
}

func makeSubbranchMarkup(subbranches map[string]string) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for key, value := range subbranches {
		button := tgbotapi.NewInlineKeyboardButtonData(value, key)
		row := tgbotapi.NewInlineKeyboardRow(button)
		rows = append(rows, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
