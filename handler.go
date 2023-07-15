package gleam

import "github.com/gempir/go-twitch-irc/v4"

type Message twitch.PrivateMessage

type Handler func(Message) Event
