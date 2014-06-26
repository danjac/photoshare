package routes

import (
	"github.com/igm/pubsub"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
)

var pub pubsub.Publisher

func sendMessage(msg string) {
	pub.Publish(msg)
}

func messageHandler(session sockjs.Session) {
	var closedSession = make(chan struct{})
	go func() {
		reader, _ := pub.SubChannel(nil)
		for {
			select {
			case <-closedSession:
				return
			case msg := <-reader:
				if err := session.Send(msg.(string)); err != nil {
					log.Println(err)
					return
				}

			}
		}
	}()
	for {
		if msg, err := session.Recv(); err == nil {
			pub.Publish(msg)
			continue
		}
		break
	}
	close(closedSession)
}
