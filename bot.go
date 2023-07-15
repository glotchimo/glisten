package gleam

import (
	"fmt"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

type BotOptions struct {
	Channel      string
	Username     string
	ClientID     string
	ClientSecret string
}

type Bot struct {
	Client *twitch.Client
	Events chan Event
	Errors chan error

	handlers map[string]Handler
	options  BotOptions
	tokens   struct {
		access  string
		refresh string
	}
}

func NewBot(opts *BotOptions) (*Bot, error) {
	var bot Bot
	bot.options = *opts

	if err := authenticate(&bot); err != nil {
		return nil, fmt.Errorf("error authenticating with Twitch: %w", err)
	}

	bot.Events = make(chan Event)
	bot.Errors = make(chan error)

	return &bot, nil
}

func (b *Bot) AddHandler(trigger string, handler Handler) {
	b.handlers[trigger] = handler
}

func (b *Bot) Connect() {
	b.Client = twitch.NewClient(b.options.Username, "oauth:"+b.tokens.access)

	b.Client.OnConnect(func() {
		fmt.Printf("bot connected to %s as %s\n", b.options.Channel, b.options.Username)
	})

	b.Client.OnPrivateMessage(func(m twitch.PrivateMessage) {
		components := strings.Split(m.Message, " ")
		if len(components) < 1 {
			return
		}
		cmd := strings.TrimSpace(components[0])

		for trigger, handler := range b.handlers {
			if cmd == trigger {
				event := handler(Message(m))
				b.Events <- event
				fmt.Println("handled event:", event.String())
				return
			}
		}
	})

	b.Client.OnNoticeMessage(func(m twitch.NoticeMessage) {
		b.Errors <- fmt.Errorf("error in notice callback: %s", m.Message)
	})

	b.Client.Join(b.options.Channel)
	if err := b.Client.Connect(); err != nil {
		b.Errors <- fmt.Errorf("error connecting to channel: %w", err)
	}
}
