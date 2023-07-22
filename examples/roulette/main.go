package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/glotchimo/gleam"
)

var bank = map[string]int{}

func main() {
	// Create a bot with credentials in the environment
	bot, err := gleam.NewBot(&gleam.BotOptions{
		Channel:  os.Getenv("TWITCH_CHANNEL"),
		Username: os.Getenv("TWITCH_USERNAME"),
		Password: os.Getenv("TWITCH_PASSWORD"),
	})
	if err != nil {
		log.Fatal("error creating new bot: ", err)
	}

	// Add a handler that fires a randomized roulette event with the given bet
	bot.AddHandler("!roulette", func(m gleam.Message) gleam.Event {
		components := strings.Split(m.Message, " ")

		var bet int
		var err error
		if len(components) < 2 {
			return gleam.Event{}
		} else if bet, err = strconv.Atoi(components[1]); err != nil {
			return gleam.Event{}
		}

		return gleam.Event{
			Type:     "roulette",
			UserID:   m.User.ID,
			Username: m.User.Name,
			Data:     bet,
		}
	})

	// Add a handler that returns the given user's points
	bot.AddHandler("!points", func(m gleam.Message) gleam.Event {
		return gleam.Event{
			Type:     "points",
			UserID:   m.User.ID,
			Username: m.User.Name,
		}
	})

	// Connect the bot
	go bot.Connect()

	// Watch for events, errors, and timers
	for {
		select {
		case event := <-bot.Events:
			switch event.Type {
			case "roulette":
				bet := event.Data.(int)
				has, ok := bank[event.UserID]
				if !ok {
					has = 0
				}

				if has == 0 {
					bet = 10
				} else if bet > has {
					bet = has
				}

				rand.Seed(int64(time.Now().Nanosecond()))
				res := rand.Intn(10)
				if res > 5 {
					bot.Say(fmt.Sprintf("You won %d points", bet*2))
					bank[event.UserID] += bet
				} else {
					bot.Say("You lost")
					bank[event.UserID] -= bet
				}
			case "points":
				has, ok := bank[event.UserID]
				if !ok {
					has = 0
				}
				bot.Say(fmt.Sprintf("You have %d points", has))
			default:
				log.Println("unsupported event received")
				continue
			}
		case error := <-bot.Errors:
			log.Fatal("got an error: ", error)
		}
	}
}
