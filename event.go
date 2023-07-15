package gleam

import "fmt"

type Event struct {
	Type     string
	UserID   string
	Username string
	Data     any
}

func (e Event) String() string {
	return fmt.Sprintf("%s from %s (%s) - %v", e.Type, e.UserID, e.Username, e.Data)
}
