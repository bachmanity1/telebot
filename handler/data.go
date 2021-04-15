package handler

import (
	"fmt"
	"strings"
	"telebot/scraper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

const maxButtonLength = 40

type updateType int

const (
	plainText = updateType(iota)
	callbackQuery
)

type reply struct {
	field    string
	text     string
	isMarkup bool
	markup   tgbotapi.InlineKeyboardMarkup
}

var replies = []reply{
	{field: "name", text: "Enter your full name (EXACTLY as it appears in your ARC)", isMarkup: false},
	{field: "branch", text: "Choose Immigration Branch", isMarkup: true, markup: branchMarkup},
	{field: "booth", text: "Choose Booth Category", isMarkup: true},
	{field: "purpose", text: "Choose purpose of visit", isMarkup: true, markup: purposeMarkup},
	{field: "phone", text: "Enter your phone number (optional)", isMarkup: false},
	{field: "receipt", text: "Receipt PlaceHolder", isMarkup: true, markup: receiptMarkup},
}

func sanitizeData(data map[string]string) {
	data["name"] = strings.ToUpper(data["name"])
	phone := getPhoneNumber(data["phone"])
	for i, val := range phone {
		key := fmt.Sprintf("phone%d", i)
		data[key] = val
	}
}

func getPhoneNumber(input string) []string {
	n := 3
	number := make([]string, 0)
	temp := make([]byte, 0)
	for i := 0; i < len(input); i++ {
		if input[i] >= '0' && input[i] <= '9' {
			temp = append(temp, input[i])
			if len(temp) == n {
				number = append(number, string(temp))
				temp = make([]byte, 0)
				n = 4
			}
		} else if input[i] == '-' {
			continue
		} else {
			return nil
		}
	}
	if len(temp) != 0 || len(number) != 3 {
		return nil
	}
	return number
}

func makeReceipt(data map[string]string) string {
	name := data["name"]
	branch := branchMap[data["branch"]]
	booth := boothMap[data["booth"]]
	date := data["resvYmd"]
	receipt := fmt.Sprintf("Successfully made an reservation!\n\n"+
		"Name: %s\nBranch: %s\nBooth: %s\nDate: %s\n\n",
		name, branch, booth, date)
	return receipt
}

var (
	boothMarkup map[string]tgbotapi.InlineKeyboardMarkup
	boothMap    map[string]string
)

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

var branchMap = map[string]string{
	"1270667": "Seoul Immigration Office",
	"1271020": "Sejongno Branch Office",
	"1270700": "Incheon Immigration Office",
	"1272143": "Incheon Immigration Office Ansan Branch Office",
	"1270947": "Suwon Immigration Office",
	"1270698": "Ulsan  Immigration Office",
	"1270686": "Busan Immigration Office",
	"1270727": "Daejeon Immigration Office",
}

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
		tgbotapi.NewInlineKeyboardButtonData("Look for an earlier date!", "continue"),
	),
)

func makeBoothesMarkup(boothes map[string][]string) {
	boothMarkup = make(map[string]tgbotapi.InlineKeyboardMarkup)
	boothMap = make(map[string]string)
	for branch, boothz := range boothes {
		boothMarkup[branch] = makeBoothMarkup(boothz)
	}
}

func makeBoothMarkup(boothz []string) tgbotapi.InlineKeyboardMarkup {
	if len(boothz)%2 != 0 {
		return tgbotapi.InlineKeyboardMarkup{}
	}
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for i := 0; i < len(boothz); i += 2 {
		key := boothz[i]
		val := boothz[i+1]
		boothMap[key] = val
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

func InitData(config *viper.Viper) error {
	boothes, err := scraper.GetBoothes(config)
	if err != nil {
		log.Debugw("InitData", "error", err)
	}
	makeBoothesMarkup(boothes)
	return nil
}
