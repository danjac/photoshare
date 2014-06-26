package routes

import (
	"encoding/json"
	"github.com/igm/pubsub"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
)

var pub pubsub.Publisher

type Message struct {
	UserName string `json:"username"`
	PhotoID  int64  `json:"photoID"`
	Type     string `json:"type"`
}

func sendMessage(msg *Message) {
	pub.Publish(msg)
}

func messageHandler(session sockjs.Session) {
	go func() {
		reader, _ := pub.SubChannel(nil)
		for {
			select {
			case msg, ok := <-reader:
				if !ok {
					log.Println("channel closed")
					return
				}
				msg = msg.(*Message)
				if body, err := json.Marshal(msg); err == nil {
					log.Println("message:", string(body))
					if err = session.Send(string(body)); err != nil {
						log.Println(err)
						return
					}
				}
			}
		}
	}()
}
