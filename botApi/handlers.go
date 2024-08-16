package main

//
//import (
//	"errors"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"net/url"
//	"strconv"
//)
//
//// HandleTelegramWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
//func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
//
//	// Parse incoming request
//	var update, err = parseTelegramRequest(r)
//	if err != nil {
//		log.Printf("error parsing update, %s", err.Error())
//		return
//	}
//
//	// Sanitize input
//	var sanitizedSeed = sanitize(update.Message.Text)
//
//	// Call RapLyrics to get a punchline
//	var lyric, errRapLyrics = getPunchline(sanitizedSeed)
//	if errRapLyrics != nil {
//		log.Printf("got error when calling RapLyrics API %s", errRapLyrics.Error())
//		return
//	}
//
//	// Send the punchline back to Telegram
//	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.Id, lyric)
//	if errTelegram != nil {
//		log.Printf("got error %s from telegram, response body is %s", errTelegram.Error(), telegramResponseBody)
//	} else {
//		log.Printf("punchline %s successfully distributed to chat id %d", lyric, update.Message.Chat.Id)
//	}
//}
//
//// parseTelegramRequest handles incoming update from the Telegram web hook
//func parseTelegramRequest(r *http.Request) (*Update, error) {
//	var update Update
//	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
//		log.Printf("could not decode incoming update %s", err.Error())
//		return nil, err
//	}
//	if update.UpdateId == 0 {
//		log.Printf("invalid update id, got update id = 0")
//		return nil, errors.New("invalid update id of 0 indicates failure to parse incoming update")
//	}
//	return &update, nil
//}
//
//// sanitize remove clutter like /start /punch or the bot name from the string s passed as input
//func sanitize(s string) string {
//	if len(s) >= lenStartCommand {
//		if s[:lenStartCommand] == startCommand {
//			s = s[lenStartCommand:]
//		}
//	}
//
//	if len(s) >= lenPunchCommand {
//		if s[:lenPunchCommand] == punchCommand {
//			s = s[lenPunchCommand:]
//		}
//	}
//	if len(s) >= lenBotTag {
//		if s[:lenBotTag] == botTag {
//			s = s[lenBotTag:]
//		}
//	}
//	return s
//}
//
//// getPunchline calls the RapLyrics API to get a punchline back.
//func getPunchline(seed string) (string, error) {
//	rapLyricsResp, err := http.PostForm(
//		rapLyricsApi,
//		url.Values{"input": {seed}})
//	if err != nil {
//		log.Printf("error while calling raplyrics %s", err.Error())
//		return "", err
//	}
//	var punchline Lyric
//	if err := json.NewDecoder(rapLyricsResp.Body).Decode(&punchline); err != nil {
//		log.Printf("could not decode incoming punchline %s", err.Error())
//		return "", err
//	}
//	defer rapLyricsResp.Body.Close()
//	return punchline.Punch, nil
//}
//
//// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
//func sendTextToTelegramChat(chatId int, text string) (string, error) {
//
//	log.Printf("Sending %s to chat_id: %d", text, chatId)
//	response, err := http.PostForm(
//		telegramApi,
//		url.Values{
//			"chat_id": {strconv.Itoa(chatId)},
//			"text":    {text},
//		})
//
//	if err != nil {
//		log.Printf("error when posting text to the chat: %s", err.Error())
//		return "", err
//	}
//	defer response.Body.Close()
//
//	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
//	if errRead != nil {
//		log.Printf("error in parsing telegram answer %s", errRead.Error())
//		return "", err
//	}
//	bodyString := string(bodyBytes)
//	log.Printf("Body of Telegram Response: %s", bodyString)
//
//	return bodyString, nil
//}
