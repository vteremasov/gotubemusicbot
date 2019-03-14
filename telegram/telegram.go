package telegram

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/rylio/ytdl"
	"github.com/vteremasov/go-music-bot/id3"
	"github.com/vteremasov/go-music-bot/storage"
	"github.com/vteremasov/go-music-bot/youtube"
	"log"
	"os"
	"reflect"
)

func Bot() {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	fmt.Println(len(updates))
	fmt.Println("Updates")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			fmt.Println("Got message:")
			fmt.Println("******************")
			fmt.Println(update.Message.Text)
			fmt.Println("******************")

			switch update.Message.Text {
			case "/start":

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, i'm a youtube music bot.")
				bot.Send(msg)

			case "/number_of_users":

				if os.Getenv("DB_SWITCH") == "on" {

					num, err := storage.GetNumberOfUsers()
					if err != nil {

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
						bot.Send(msg)
					}

					ans := fmt.Sprintf("%d people used me for listenning music from youtube", num)

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, ans)
					bot.Send(msg)
				} else {

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database not connected, so i can't say you how many peoples used me.")
					bot.Send(msg)
				}
			default:
				fmt.Println("Trying to work")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Started working on your request.")
				_, err = bot.Send(msg)

				if err != nil {
					log.Panic(err)
				}

				if youtube.IsYoutubeLink(update.Message.Text) {
					sendSong(update, bot)
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I don't know what does it mean: "+update.Message.Text)
					_, err = bot.Send(msg)

					if err != nil {
						log.Panic(err)
					}
				}
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry I don't understand that...")
			bot.Send(msg)
		}
	}
}

func sendSong(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	fmt.Println(update.Message.Text)
	youtubeID := youtube.GetYoutubeID(update.Message.Text)
	videoInfo, _ := ytdl.GetVideoInfoFromID(youtubeID)
	mp3Location := youtube.DownloadMp3(*videoInfo, "/var/tmp")
	meta := id3.CollectMetadata(videoInfo)
	finalLocation := youtube.RenameMp3File(mp3Location, meta)
	id3.SetMetadata(finalLocation, meta)
	name := tgbotapi.NewMessage(update.Message.Chat.ID, videoInfo.Title+" is ready.")
	_, _ = bot.Send(name)
	file, err := os.Open(finalLocation)
	fileStat, _ := file.Stat()
	fileReader := tgbotapi.FileReader{
		Reader: file,
		Name:   meta.Title,
		Size:   fileStat.Size(),
	}
	if err != nil {
		log.Panic(err)
	}
	// TODO: Save file info to db
	fmt.Println(name)
	if os.Getenv("DB_SWITCH") == "on" {

		if err := storage.CollectData(update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.Text); err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error, but bot still working.")
			bot.Send(msg)
		}
	}
	audioMsg := tgbotapi.NewAudioUpload(update.Message.Chat.ID, fileReader)
	_, err = bot.Send(audioMsg)
	if err != nil {
		log.Panic(err)
	}
}
