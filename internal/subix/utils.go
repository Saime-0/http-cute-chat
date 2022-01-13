package subix

import (
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"strings"
)

func (s *Subix) writeToUsers(ids []int, body model.EventResult) {

	eventType := getEventType(body)
	for _, id := range ids {
		user, ok := s.users[id]
		if !ok {
			continue
		}

		next := user.Root.next
		for next != nil {
			select {
			case next.Ch <- &model.SubscriptionBody{
				//Rev:   next.rev[eventType][objectID],
				Event: eventType,
				Body:  body,
			}:
				//next.rev[eventType] += 1 todo fix
				fmt.Printf("Message write to client chan %p (id:%d)\n", next, id) // debug
			default:
				fmt.Printf("client chan %p (id:%d) is close.. skip\n", next, id) // debug
				delete(s.users, id)
			}
			next = next.next
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
