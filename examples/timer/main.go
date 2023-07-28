package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glotchimo/glisten"
	"github.com/glotchimo/glisten/handlers"
)

type timer struct {
	Username string
	Duration time.Duration
}

func (t *timer) Start(ch chan *timer) {
	timer := time.NewTimer(t.Duration)
	go func() {
		<-timer.C
		ch <- t
	}()
}

func main() {
	// Create a bot with credentials in the environment
	bot, err := glisten.NewBot(&glisten.BotOptions{
		Channel:  os.Getenv("TWITCH_CHANNEL"),
		Username: os.Getenv("TWITCH_USERNAME"),
		Password: os.Getenv("TWITCH_PASSWORD"),
	})
	if err != nil {
		log.Fatal("error creating new bot: ", err)
	}

	// Add a handler that creates an event to start a timer
	bot.AddHandler("!timer", handlers.Timer)

	// Connect the bot
	go bot.Connect()

	// Watch for events, errors, and timers
	timers := make(chan *timer)
	for {
		select {
		case event := <-bot.Events:
			switch event.Type {
			case "timer":
				duration := event.Data.(time.Duration)
				timer := &timer{
					Username: event.Username,
					Duration: duration,
				}
				timer.Start(timers)
				msg := fmt.Sprintf("%s's timer is set (%s)", event.Username, duration.String())
				bot.Say(msg)
			}

		case error := <-bot.Errors:
			log.Fatal("got an error: ", error)

		case timer := <-timers:
			msg := fmt.Sprintf("%s's timer is up", timer.Username)
			bot.Say(msg)
		}
	}
}
