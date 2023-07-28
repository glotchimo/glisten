package handlers

import (
	"strconv"
	"strings"

	"github.com/glotchimo/glisten"
)

func Roulette(m glisten.Message) glisten.Event {
	components := strings.Split(m.Message, " ")

	var bet int
	var err error
	if len(components) < 2 {
		return glisten.Event{}
	} else if bet, err = strconv.Atoi(components[1]); err != nil {
		return glisten.Event{}
	}

	return glisten.Event{
		Type:     "roulette",
		UserID:   m.User.ID,
		Username: m.User.Name,
		Data:     bet,
	}
}
