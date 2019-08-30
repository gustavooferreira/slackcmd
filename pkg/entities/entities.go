package entities

type RequestContext struct {
	Cmd         string
	Text        string
	ResponseURL string
	TriggerID   string
	TeamID      string
	ChannelID   string
	UserID      string
}
