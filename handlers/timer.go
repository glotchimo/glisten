package handlers

import (
	"strings"
	"time"

	"github.com/glotchimo/glisten"
)

func Timer(m glisten.Message) glisten.Event {
	components := strings.Split(m.Message, " ")

	var duration time.Duration
	var err error
	if len(components) < 2 {
		return glisten.Event{}
	} else if duration, err = time.ParseDuration(components[1]); err != nil {
		return glisten.Event{}
	}

	return glisten.Event{
		Type:     "timer",
		UserID:   m.User.ID,
		Username: m.User.Name,
		Data:     duration,
	}
}
