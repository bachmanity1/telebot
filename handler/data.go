package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

const maxButtonLength = 40

var branchMarkup = tgbotapi.NewInlineKeyboardMarkup(
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

var purposeMarkup = tgbotapi.NewInlineKeyboardMarkup(
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

var receiptMarkup = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Good Enough!", "exit"),
	), tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Look for an earlier date", "2"),
	),
)

func makeBoothMarkup(boothes []string) tgbotapi.InlineKeyboardMarkup {
	if len(boothes)%2 != 0 {
		return tgbotapi.InlineKeyboardMarkup{}
	}
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for i := 0; i < len(boothes); i += 2 {
		key := boothes[i]
		val := boothes[i+1]
		button := tgbotapi.NewInlineKeyboardButtonData(shorten(val), key)
		row := tgbotapi.NewInlineKeyboardRow(button)
		rows = append(rows, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func shorten(line string) string {
	if len(line) > maxButtonLength {
		n := len(line) - maxButtonLength
		line = "..." + line[n:]
	}
	return line
}

type reply struct {
	field    string
	text     string
	isMarkup bool
	markup   tgbotapi.InlineKeyboardMarkup
}

var replies = []reply{
	{field: "username", text: "Enter your username", isMarkup: false},
	{field: "password", text: "Enter your password", isMarkup: false},
	{field: "branch", text: "Choose Immigration Branch", isMarkup: true, markup: branchMarkup},
	{field: "booth", text: "Choose Booth Category", isMarkup: true},
	{field: "purpose", text: "Choose purpose of visit", isMarkup: true, markup: purposeMarkup},
	{field: "phone", text: "Enter your phone number", isMarkup: false},
	{field: "receipt", text: "Receeipt", isMarkup: true, markup: receiptMarkup},
}
