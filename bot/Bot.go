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
		tgbotapi.NewKeyboardButton("Rules"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Russian"),
		tgbotapi.NewKeyboardButton("English"),
	))

type stringSlice []string

func main() {
	bot, err := tgbotapi.NewBotAPI("") // just test bot api
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
	userCities := make(map[int]stringSlice)
	lastCity := make(map[int]string)
	used := make(map[int]stringSlice)
	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		uID := update.Message.From.ID
		//Обработка команд keyboard
		if update.Message.Text == "English" {
			msg.Text = "Not ready yet"
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
			continue
		} else if update.Message.Text == "Russian" {
			msg.Text = "Выбран русский язык"
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
			continue
		} else if update.Message.Text == "End game" {
			used[uID] = used[uID][:0]
			log.Println(used)
			lastCity[uID] = ""
			msg.Text = "Игра сброшена \nСпасибо за игру"
			msg.ReplyToMessageID = update.Message.MessageID
			userCities[uID] = cities
			_, _ = bot.Send(msg)
			continue
		} else if update.Message.Text == "Start game" {
			used[uID] = used[uID][:0]
			lastCity[uID] = ""
			msg.Text = "Для начала игры отправьте название города"
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
			userCities[uID] = cities
			continue
		} else if update.Message.Text == "Rules" {
			msg.Text = "Отправьте название города. Бот пришлет вам следующий город, начинающийся с той же буквы. Ответьте ему тем же"
			_, _ = bot.Send(msg)
			continue
		} else if update.Message.Text == "/start" {
			continue
		} else if update.Message.Text == "/start_game" {
			msg.ReplyMarkup = keyboard
			msg.Text = "Начните игру, нажав на Start game"
			_, _ = bot.Send(msg)
			continue
		} else if update.Message.Text == "/help" {
			msg.Text = "Начните игру, написав название города. Бот отправит вам город, название которого начинается  на последнюю букву вашего города. Продолжайте игру"
			_, _ = bot.Send(msg)
			continue
		}
		//до сюда
		if len(userCities[uID]) == 0 {
			msg.Text = "Вы не начали игру. Нажмите /start_game или \"Start game\""
			msg.ReplyMarkup = keyboard
			_, _ = bot.Send(msg)
			continue
		}
		current := update.Message.Text
		log.Printf("[%s]: %s", update.Message.From.UserName, strings.ToLower(current))
		var toSend string
		toSend, userCities[uID], lastCity[uID] = chooseWord(current, userCities[uID], lastCity[uID], used[uID])
		used[uID] = append(used[uID], lastCity[uID], current)

		msg.Text = toSend
		if update.Message.Text == "English" {
			msg.Text = "Not ready yet"
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
			continue
		} else if update.Message.Text == "Russian" {
			msg.Text = "Выбран русский язык"
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
			continue
		}
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
	}
}

func chooseWord(current string, cities []string, lastSent string, used stringSlice) (string, []string, string) {
	result := " "
	if lastSent != "" && []rune(lastSent)[len([]rune(lastSent))-1] != []rune(strings.ToLower(current))[0] {
		result = "Ваше слово не подходит. Последний город - " + lastSent
		return result, cities, lastSent
	}
	f := stringSlice(cities).contains(current)
	if f > -1 {
		cities = append(cities[:f], cities[f+1:]...)
	} else if used.contains(current) != -1 {
		result = "Этот город уже был сыгран. Попробуйте еще"
		return result, cities, lastSent
	} else {
		result = "Такого города не существует[вероятно, в нашей бд]. Попробуйте еще"
		return result, cities, lastSent
	}
	if []rune(current)[len([]rune(current))-1] == rune('ь'){
		current = string([]rune(current)[0:len([]rune(current))-1])
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
