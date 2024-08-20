package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	// Menu texts
	firstMenu  = "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	secondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."

	// Button texts
	nextButton     = "Next"
	backButton     = "Back"
	tutorialButton = "Tutorial"

	// Store bot screaming status
	screaming = false
	bot       *tgbotapi.BotAPI

	// Keyboard layout for the first menu. One button, one row
	firstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(nextButton, nextButton),
		),
	)

	// Keyboard layout for the second menu. Two buttons, one per row
	secondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(tutorialButton, "https://core.telegram.org/bots/api"),
		),
	)
)

func init() {
	mode := os.Getenv(gin.EnvGinMode)
	gin.SetMode(mode)
}

func main() {
	bot, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Panic(err)
	}

	router := gin.Default()

	enverr := godotenv.Load()
	if enverr != nil {
		log.Fatal("Error loading .env file")
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}

	router.GET("/status", getServerStatus)
	router.Run(":8080")
}

//func main() {
//	var err error
//	router := gin.Default()
//
//	enverr := godotenv.Load()
//	if enverr != nil {
//		log.Fatal("Error loading .env file")
//	}
//
//	bot, err = tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
//	if err != nil {
//		// Abort if something is wrong
//		log.Panic(err)
//	}
//	log.Printf("Authorized on account %s", bot.Self.UserName)
//
//	// Set this to true to log all interactions with telegram servers
//	bot.Debug = false
//
//	u := tgbotapi.NewUpdate(0)
//	u.Timeout = 60
//
//	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
//	ctx := context.Background()
//	ctx, cancel := context.WithCancel(ctx)
//
//	// `updates` is a golang channel which receives telegram updates
//	updates := bot.GetUpdatesChan(u)
//
//	// Pass cancellable context to goroutine
//	go receiveUpdates(ctx, updates)
//
//	// Tell the user the bot is online
//	log.Println("Start listening for updates. Press enter to stop")
//
//	// Wait for a newline symbol, then cancel handling updates
//	bufio.NewReader(os.Stdin).ReadBytes('\n')
//	cancel()
//
//	router.GET("/status", getServerStatus)
//	router.Run(":8080")
//}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	log.Printf("Listening for updates %v", updates)
	for {
		log.Println("loop")
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			log.Println("Stopping updates")
			return
		// receive update from channel and then handle it
		case update := <-updates:
			log.Println("Received update")
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	log.Printf("Received update: %v", update)
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message)
		break

	// Handle button clicks
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery)
		break
	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text
	log.Printf("Received message from %s: %s", user.FirstName, text)
	if user == nil {
		return
	}

	// Print to console
	log.Printf("%s wrote %s", user.FirstName, text)

	var err error
	if strings.HasPrefix(text, "/") {
		err = handleCommand(message.Chat.ID, text)
	} else if screaming && len(text) > 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, strings.ToUpper(text))
		// To preserve markdown, we attach entities (bold, italic..)
		msg.Entities = message.Entities
		_, err = bot.Send(msg)
	} else {
		// This is equivalent to forwarding, without the sender's name
		copyMsg := tgbotapi.NewCopyMessage(message.Chat.ID, message.Chat.ID, message.MessageID)
		_, err = bot.CopyMessage(copyMsg)
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func handleCommand(chatId int64, command string) error {
	var err error

	switch command {
	case "/scream":
		screaming = true
		break

	case "/whisper":
		screaming = false
		break

	case "/menu":
		err = sendMenu(chatId)
		break
	}

	return err
}

func handleButton(query *tgbotapi.CallbackQuery) {
	var text string

	markup := tgbotapi.NewInlineKeyboardMarkup()
	message := query.Message

	if query.Data == nextButton {
		text = secondMenu
		markup = secondMenuMarkup
	} else if query.Data == backButton {
		text = firstMenu
		markup = firstMenuMarkup
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	bot.Send(callbackCfg)

	// Replace menu text and keyboard
	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	bot.Send(msg)
}

func sendMenu(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, firstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	_, err := bot.Send(msg)
	return err
}

type Status struct {
	Status string `json:"status"`
}

var status Status = Status{Status: "ok"}

func getServerStatus(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, status)
}
