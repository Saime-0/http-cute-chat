package subix

import (
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"strings"
)

func (s *Subix) writeToMembers(membersID []int, body model.EventResult) {
	eventType := getEventType(body)
	for _, memberID := range membersID {
		member, ok := s.members[memberID]
		if !ok {
			continue
		}
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

func (s *Subix) writeToChats(chatsID []int, body model.EventResult) {
	eventType := getEventType(body)
	for _, chatID := range chatsID {
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

func (s *Subix) writeToUsers(usersID []int, body model.EventResult) {
	eventType := getEventType(body)
	for _, userID := range usersID {
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
		//fmt.Printf("client %p (id:%d) marked.. skip\n", client, (*client).UserID)
		return
	}
	select {
	case (*client).Ch <- subbody:
		//fmt.Printf("Message write to client %s (id:%d)\n", client.sessionKey, (*client).UserID)

	default:
		//fmt.Printf("client chan %p (id:%d) is close.. skip\n", client, (*client).UserID)
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
