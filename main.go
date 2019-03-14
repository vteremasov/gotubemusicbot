package main

import (
	"fmt"
	"github.com/vteremasov/go-music-bot/storage"
	"github.com/vteremasov/go-music-bot/telegram"
	"os"
	"time"
)

func main() {
	fmt.Println("Starting...")
	time.Sleep(10 * time.Second)

	if os.Getenv("CREATE_TABLE") == "yes" {

		if os.Getenv("DB_SWITCH") == "on" {

			if err := storage.CreateTable(); err != nil {

				panic(err)
			}
		}
	}

	time.Sleep(10 * time.Second)

	fmt.Println("Started...")
	telegram.Bot()
}
