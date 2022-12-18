package pkg

import (
	"database/sql"
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
		// Send an update for the new item to the Telegram channel
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
			return err
		}
	}
	return nil
}
