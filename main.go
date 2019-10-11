/* TODO
 * - error handling
 * - check on start if all env are set
 * - check for chat token so that only the correct person gets the notifications
 * - chat token could be set in .env or passed in via command line 
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

	"github.com/robfig/cron"
)

type GithubDataType struct {
	Subject Subject `json:"subject"`
	Repository Repository `json:"repository"`
}

type Subject struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type Repository struct {
	Name string `json:"name"`
}

type GithubDataMessage struct {
	GithubData GithubDataType
	Send bool
}

var telegram_api_token string
var github_api_token string
var c = cron.New()
var dataSet []GithubDataMessage

func init() {
	gotenv.Load()
	telegram_api_token = os.Getenv("TELEGRAM_API_TOKEN") 
	github_api_token = os.Getenv("GITHUB_API_TOKEN")
}

func main() {
	updater, err := gotgbot.NewUpdater(telegram_api_token)

	if err != nil {
		log.Fatal(err)
	}

	// message handler
	updater.Dispatcher.AddHandler(handlers.NewCommand("start", addCronJob))

	// start getting updates
	updater.StartPolling()

	// wait
	updater.Idle()
}

func addIfNotIncluded(myslice []GithubDataMessage, item GithubDataMessage) []GithubDataMessage {
	for i := 0; i < len(myslice); i++ {
		if myslice[i].GithubData == item.GithubData {
			return myslice
		}
	}
	return append(myslice, item)
}
  
func addCronJob(bot ext.Bot, update *gotgbot.Update) error {
	if len(c.Entries()) < 1 {
		// call function every minute to check for updates
		c.AddFunc("* * * * *", func() {
			checkGithub(bot, update)
		})
		c.Start()
	}

	return nil
}

func checkGithub(bot ext.Bot, update *gotgbot.Update) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	request, err := http.NewRequest("GET", "https://api.github.com/notifications", nil)
	request.Header.Add("Authorization", "token " + github_api_token)

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

	var githubData []GithubDataType

	error := json.Unmarshal(body, &githubData)

	if error != nil {
		log.Fatal(error)
	}

	// add the the github data if not already included
	for i := 0; i < len(githubData); i++ {
		dataSet = addIfNotIncluded(dataSet, GithubDataMessage{githubData[i], false})
	}

	for index, value := range dataSet {
		if dataSet[index].Send == false {
			_, err := bot.SendMessage(update.Message.Chat.Id, "Repository: " + value.GithubData.Repository.Name + "\nNotification: " + value.GithubData.Subject.Title + "\n" + value.GithubData.Subject.Url)
			dataSet[index].Send = true
			if err != nil {
				log.Fatal(err)
			}

		}
	}

	return nil
}
