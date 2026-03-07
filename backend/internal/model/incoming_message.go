package model

type IncomingMessage struct {
	From string `json:"from"`
	Body string `json:"body"`
}