package main

import (
	"bufio"
	"encoding/csv"
	"github.com/Syfaro/telegram-bot-api"
	"io"
	"log"
	"os"
	"strings"
)

var keyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Start game"),
		tgbotapi.NewKeyboardButton("End game"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Russian"),
		tgbotapi.NewKeyboardButton("English"),
	))

type stringSlice []string

func main() {
	bot, err := tgbotapi.NewBotAPI("945600369:AAHVNGwrXhbT1KIAa6y5LV5zrC1gAiXVgRs") // just test bot api
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	csvFile, err := os.Open("bot/city.csv")
	defer csvFile.Close()
	if err != nil {
		log.Panic(err)
	}
	csvReader := csv.NewReader(bufio.NewReader(csvFile))

	var cities []string
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		cities = append(cities, string(line[3]))
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	var lastCityM map[int]string
	lastCityM = make(map[int]string)
	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		current := update.Message.Text
		log.Printf("[%s]: %s", update.Message.From.UserName, strings.ToLower(current))
		var toSend string
		toSend, cities, lastCityM[update.Message.From.ID] = chooseWord(current, cities, lastCityM[update.Message.From.ID])

		msg.Text = toSend
		if update.Message.Text == "English" {
			msg.Text = "Not ready yet"
		} else if update.Message.Text == "Russian" {
			msg.Text = "Выбран русский язык"
		}
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
	}
}

func chooseWord(current string, cities []string, lastSent string) (string, []string, string) {
	result := " "
	if lastSent != "" && []rune(lastSent)[len([]rune(lastSent))-1] != []rune(strings.ToLower(current))[0] {
		result = "Ваше слово не подходит. Последний город - " + lastSent
		return result, cities, lastSent
	}
	f := stringSlice(cities).contains(current)
	if f > -1 {
		cities = append(cities[:f], cities[f+1:]...)
	} else {
		result = "Название уже было использовано или такого города не существует[вероятно, в нашей бд]. Попробуйте еще"
		return result, cities, lastSent
	}
	for i, city := range cities {
		if []rune(strings.ToLower(city))[0] == []rune(strings.ToLower(current))[len([]rune(current))-1] {
			result = city
			cities = append(cities[:i], cities[i+1:]...)
			return result, cities, result
		}
	}
	if result == " " {
		result = "Это невозможно, но вы выиграли"
		return result, cities, lastSent
	}
	return result, cities, lastSent
}

func (s stringSlice) contains(e string) int {
	for i, a := range s {
		if strings.ToLower(a) == strings.ToLower(e) {
			return i
		}
	}
	return -1
}
