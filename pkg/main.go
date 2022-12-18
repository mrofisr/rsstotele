package pkg

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
)

func CheckForUpdates(bot *tgbotapi.BotAPI, channelName string, db *sql.DB) error {
	// Parse the RSS feed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(os.Getenv("RSS_URL"))
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Get the 5 newest items in the RSS feed
	var newestItems []*gofeed.Item
	for i := 0; i < 5 && i < len(feed.Items); i++ {
		newestItems = append(newestItems, feed.Items[i])
	}
	// Send updates for the newest items to the Telegram channel
	for _, item := range newestItems {
		// Check if the item has already been sent
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM items WHERE link = ?", item.Link).Scan(&count)
		if err != nil {
			log.Fatal(err)
		}

		// If the item has not been sent, send an update for it to the Telegram channel
		if count == 0 {
			fmt.Printf("Sending: %s\n", item.Title)
			msg := tgbotapi.NewMessageToChannel(channelName, item.Title+"\n"+item.Link)
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second)

			// Insert the new item into the database
			_, err = db.Exec("INSERT INTO items (title, link, published_at) VALUES (?, ?, ?)", item.Title, item.Link, item.Published)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return nil
}
