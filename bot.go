package glisten

import (
	"fmt"
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

// BotOptions are the values needed for connecting the bot to a channel/user
type BotOptions struct {
	Channel  string
	Username string
	Password string
}

// Bot is an abstraction over Twitch IRC that facilitates the creation of
// command-driven chat bots.
type Bot struct {
	// Events is the channel where all events arrive after being handled by the
	// registered handlers.
	Events chan Event

	// Errors is the channel were all connection and post-connection errors are
	// sent for handling by the user.
	Errors chan error

	client   *twitch.Client
	handlers map[string]Handler
	options  BotOptions
}

// NewBot sets up a Bot with the provided options and opens its channels.
func NewBot(opts *BotOptions) (*Bot, error) {
	var bot Bot

	bot.handlers = make(map[string]Handler)
	bot.options = *opts

	bot.Events = make(chan Event)
	bot.Errors = make(chan error)

	return &bot, nil
}

// Add handler registers an event handler function on the bot's handler map.
func (b *Bot) AddHandler(trigger string, handler Handler) {
	b.handlers[trigger] = handler
}

// Connect launches the OAuth flow and, if completed, connects to Twitch IRC
// and starts listening for messages, and should be launched in a goroutine.
func (b *Bot) Connect() {
	b.client = twitch.NewClient(b.options.Username, b.options.Password)

	b.client.OnConnect(func() {
		log.Printf("bot connected to %s as %s\n", b.options.Channel, b.options.Username)
	})

	b.client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		components := strings.Split(msg.Message, " ")
		if len(components) < 1 {
			return
		}
		cmd := strings.TrimSpace(components[0])

		for trigger, handler := range b.handlers {
			if cmd == trigger {
				event := handler(Message(msg))
				b.Events <- event
				log.Println("handled event:", event.String())
				return
			}
		}
	})

	b.client.OnNoticeMessage(func(msg twitch.NoticeMessage) {
		b.Errors <- fmt.Errorf("error in notice callback: %s", msg.Message)
	})

	b.client.Join(b.options.Channel)
	if err := b.client.Connect(); err != nil {
		b.Errors <- fmt.Errorf("error connecting to channel: %w", err)
	}
}

func (b Bot) Say(msg string) {
	b.client.Say(b.options.Channel, msg)
}
