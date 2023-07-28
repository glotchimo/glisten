package glisten

import "fmt"

// Event abstracts a single event produced by a Handler after its defined
// trigger is detected by the IRC listener.
type Event struct {
	// Type is an arbitrary string to describe the type of event.
	Type string

	// UserID is the ID of the user that issued the event.
	UserID string

	// Username is the username of the user that issued the event.
	Username string

	// Data is any additional data that needs to be attached to the event.
	Data any
}

// String prints a human-friendly representation of the event.
func (e Event) String() string {
	return fmt.Sprintf("%s from %s (%s) - %v", e.Type, e.UserID, e.Username, e.Data)
}
