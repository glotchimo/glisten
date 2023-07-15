package gleam

import "github.com/gempir/go-twitch-irc/v4"

// Message is a facade type over twitch.PrivateMessage so that users don't
// need to explicitly import the go-twitch-irc package when adding handlers.
type Message twitch.PrivateMessage

// Handler is a function that takes an IRC message and produces an event.
type Handler func(Message) Event
