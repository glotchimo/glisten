package gleam

type Event struct {
	Type     string
	UserID   string
	Username string
	Data     any
}
