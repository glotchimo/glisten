package reactors

import (
	"fmt"
	"time"

	"github.com/glotchimo/glisten"
)

type Timers interface {
	Start(chan *Timer)
}

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

func TimerReactor(bot *gleam.Bot, event gleam.Event, timers chan *timer) {
	duration := event.Data.(time.Duration)
	timer := &{
		Username: event.Username,
		Duration: duration,
	}

	bot.Say(fmt.Sprintf(
		"%s's timer is set (%s)",
		event.Username,
		duration.String()))
}
