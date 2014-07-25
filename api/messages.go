package api

import (
	"encoding/json"
	"github.com/igm/pubsub"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
)

var pub pubsub.Publisher

// SocketMessage represents info to be sent
type SocketMessage struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	PhotoID  int64  `json:"photoID"`
	Type     string `json:"type"`
}

func sendMessage(msg *SocketMessage) {
	pub.Publish(msg)
}

func receiveMessage(session sockjs.Session) {
	reader, _ := pub.SubChannel(nil)
	for {
		select {
		case msg, ok := <-reader:
			if !ok {
				log.Println("channel closed")
				return
			}
			msg = msg.(*SocketMessage)
			if body, err := json.Marshal(msg); err == nil {
				log.Println("message:", string(body))
				if err = session.Send(string(body)); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

var messageHandler = sockjs.NewHandler(
	"/api/messages",
	sockjs.DefaultOptions, func(session sockjs.Session) {
		go func() {
			receiveMessage(session)
		}()
	})
