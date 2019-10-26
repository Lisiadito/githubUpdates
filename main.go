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
	Subject    Subject    `json:"subject"`
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
	Send       bool
}

var telegram_api_token string
var github_api_token string
var c = cron.New()
var dataSet []GithubDataMessage

func init() {
	log.Println("Initialization starting...")

	gotenv.Load()
	telegram_api_token = os.Getenv("TELEGRAM_API_TOKEN")
	github_api_token = os.Getenv("GITHUB_API_TOKEN")

	if(len(telegram_api_token) == 0 || len(github_api_token) == 0) {
		log.Fatal("env missing")
	}
	
	log.Println("Initialization done.")
}

func main() {
	log.Println("Bot started.")

	updater, err := gotgbot.NewUpdater(telegram_api_token)

	if err != nil {
		log.Fatal(err)
	}

	// message handler
	updater.Dispatcher.AddHandler(handlers.NewCommand("start", addCronJob))
	updater.Dispatcher.AddHandler(handlers.NewCommand("running", isRunning))

	// start getting updates
	updater.StartPolling()

	// wait
	updater.Idle()
}

func addIfNotIncluded(item GithubDataMessage) []GithubDataMessage {
	for i := 0; i < len(dataSet); i++ {
		if dataSet[i].GithubData == item.GithubData {
			return dataSet
		}
	}
	return append(dataSet, item)
}

// checks if any messages in the dataSet are not in the newMessages
// which means they are already read
func removeIfRead(newMessages []GithubDataType) {
	var tmp []GithubDataMessage
	for i := 0; i < len(dataSet); i++ {
		for j := 0; j < len(newMessages); j++ {
			if dataSet[i].GithubData == newMessages[j] {
				tmp = append(tmp, dataSet[i])
			}
		}
	}
	dataSet = tmp
}

func addCronJob(bot ext.Bot, update *gotgbot.Update) error {
	log.Println("Got start command.")
	if len(c.Entries()) < 1 {
		// call function every minute to check for updates
		c.AddFunc("* * * * *", func() {
			checkGithub(bot, update)
		})
		c.Start()
		log.Println("Added cronjob.")
	}

	return nil
}

func isRunning(bot ext.Bot, update *gotgbot.Update) error {
	bot.SendMessage(update.Message.Chat.Id, "Bot is running")
	return nil
}

func checkGithub(bot ext.Bot, update *gotgbot.Update) error {
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

	var githubData []GithubDataType

	error := json.Unmarshal(body, &githubData)

	if error != nil {
		log.Fatal(error)
	}

	// add the the github data if not already included
	for i := 0; i < len(githubData); i++ {
		dataSet = addIfNotIncluded(GithubDataMessage{githubData[i], false})
	}

	// IMPORTANT this needs to be called after addIfNotIncluded
	removeIfRead(githubData)

	for index, value := range dataSet {
		if dataSet[index].Send == false {
			_, err := bot.SendMessage(update.Message.Chat.Id, "Repository: "+value.GithubData.Repository.Name+"\nNotification: "+value.GithubData.Subject.Title+"\n"+value.GithubData.Subject.Url)
			dataSet[index].Send = true
			if err != nil {
				log.Fatal(err)
			}

		}
	}

	return nil
}
