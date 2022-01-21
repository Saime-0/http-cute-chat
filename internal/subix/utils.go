package subix

import (
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"strings"
)

func (s *Subix) writeToChats(chats []int, body model.EventResult) {
	eventType := getEventType(body)
	for _, chatID := range chats {
		chat, ok := s.chats[chatID]
		if !ok {
			continue
		}

		for _, member := range chat.members {

			for _, client := range member.clients {
				s.writeToClient(
					client,
					&model.SubscriptionBody{
						Event: eventType,
						Body:  body,
					},
				)
			}

		}

	}
}

func (s *Subix) writeToUsers(users []int, body model.EventResult) {
	eventType := getEventType(body)
	for _, userID := range users {
		user, ok := s.users[userID]
		if !ok {
			continue
		}
		for _, client := range user.clients {
			s.writeToClient(
				client,
				&model.SubscriptionBody{
					Event: eventType,
					Body:  body,
				},
			)
		}

	}
}

func (s *Subix) writeToClient(client *Client, subbody *model.SubscriptionBody) {
	if (*client).marked {
		fmt.Printf("client %p (id:%d) marked.. skip\n", client, (*client).UserID) // debug
		return
	}
	select {
	case (*client).Ch <- subbody:
		fmt.Printf("Message write to client %s (id:%d)\n", client.sessionKey, (*client).UserID) // debug

	default:
		fmt.Printf("client chan %p (id:%d) is close.. skip\n", client, (*client).UserID) // debug
		if client != nil {
			s.deleteClient(client.sessionKey)
		}
	}
}

func getEventType(body model.EventResult) string {
	bodyType := fmt.Sprintf("%T", body)
	dot := strings.LastIndex(
		bodyType,
		".",
	)
	if dot == -1 {
		panic("invalid index")
	}
	return strings.ToUpper(bodyType[dot+1:])
}
