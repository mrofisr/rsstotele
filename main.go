package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"rsstogo/pkg"
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	// Open the SQLite database
	db, err := sql.Open("sqlite3", "rss.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a table to store the RSS feed items
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id INTEGER PRIMARY KEY,
			title TEXT,
			link TEXT,
			published_at TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	// Replace YOUR_BOT_TOKEN with the API token for your bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	// Replace CHANNEL_NAME with the name of your Telegram channel (e.g. "@mychannel")
	channelName := os.Getenv("CHANNEL_USERNAME")
	// Send updates for the newest items to the Telegram channel
	s := gocron.NewScheduler(time.FixedZone("UTC+7", 7*60*60))
	fmt.Println("Start send to your telegram")
	s.Every(1).Minutes().Do(func() {
		err = pkg.CheckForUpdates(bot, channelName, db)
		if err != nil {
			log.Fatal(err)
		}
	})
	s.StartAsync()
	s.StartBlocking()
}
