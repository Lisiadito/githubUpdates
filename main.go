/* TODO
 * - error handling
 * - check on start if all env are set
 * - check for chat token so that only the correct person gets the notifications
 * - chat token could be set in .env or passed in via command line 
 * - track which updates are already sent
 * - check in intervals for updates. maybe use (https://godoc.org/github.com/robfig/cron)
 * - send message on new updates
 * - define struct and maybe send more information like repo name
 * - dockerize
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"

	"github.com/subosito/gotenv"
)

type githubDataType struct {
	Subject Subject `json:"subject"`
}

type Subject struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func init() {
	gotenv.Load()
}

func main() {

	telegram_api_token := os.Getenv("TELEGRAM_API_TOKEN")

	updater, err := gotgbot.NewUpdater(telegram_api_token)

	if err != nil {
		log.Fatal(err)
	}

	// message handler
	updater.Dispatcher.AddHandler(handlers.NewCommand("start", foo))

	// start getting updates
	updater.StartPolling()

	// wait
	updater.Idle()
}

func foo(bot ext.Bot, update *gotgbot.Update) error {
	github_api_token := os.Getenv("GITHUB_API_TOKEN")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	request, err := http.NewRequest("GET", "https://api.github.com/notifications", nil)
	request.Header.Add("Authorization", "token "+github_api_token)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))

	var githubData []githubDataType

	error := json.Unmarshal(body, &githubData)

	if error != nil {
		log.Fatal(error)
	}

	for i := 0; i < len(githubData); i++ {
		_, err := bot.SendMessage(update.Message.Chat.Id, githubData[i].Subject.Title)

		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
