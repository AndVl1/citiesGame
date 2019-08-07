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
	), )

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
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true

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

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		current := update.Message.Text
		log.Printf("[%s]: %s", update.Message.From.UserName, strings.ToLower(current))
		log.Println(len([]rune(current)), strings.ToLower(current),strings.ToLower(current)[0], strings.ToLower(current)[len(current)-1])
		toSend := " "
		i, city := 0, ""
		for i, city = range cities{
			//log.Println((city)[0], strings.ToLower(current)[len(current) - 1])
			if []rune(strings.ToLower(city))[0] == []rune(strings.ToLower(current))[len([]rune(current)) - 1] {
				toSend = city
				cities = append(cities[:i], cities[i+1:]...)
				break
			}
		}
		if toSend == " " {
			toSend = "Это невозможно, но вы выиграли"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, toSend)

		if update.Message.Text == "English" {
			msg.Text = "Not ready yet"
		}

		_, _ = bot.Send(msg)
	}
}