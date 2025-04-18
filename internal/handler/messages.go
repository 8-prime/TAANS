package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"taans/internal/telegram"
)

type Message struct {
	Text string `json:"text"`
}

func HandleNewMessage(message chan telegram.Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var m Message
		err := decoder.Decode(&m)
		if err != nil {
			log.Printf("Error decoding message: %v", err)
			return
		}
		message <- telegram.Message{Text: m.Text}
		w.WriteHeader(http.StatusOK)
	}
}
