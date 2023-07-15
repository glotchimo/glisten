package gleam

import "github.com/gempir/go-twitch-irc/v4"

type Handler func(twitch.PrivateMessage) Event
