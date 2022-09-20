package models

type Feed struct {
	Channels []Channel `json:"channels"`
}
