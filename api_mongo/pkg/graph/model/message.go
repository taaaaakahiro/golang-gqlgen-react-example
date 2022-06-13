package model

type Message struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	// see https://gqlgen.com/getting-started/#dont-eagerly-fetch-the-user
	// MEMO Userを指定されたときだけ取ってくるようにする、
	UserID string `json:"user_id"`
}
