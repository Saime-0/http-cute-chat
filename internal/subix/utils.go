package subix

import (
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (s *Subscription) writeToUsers(ids []int, body model.EventResult) {

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

func getEventType(body model.EventResult) model.EventType {
	switch body.(type) {

	case *model.NewMessage:
		return model.EventTypeNewmessage
	case *model.UpdateUser:
		return model.EventTypeUpdateunit
	case *model.UpdateMember:
		return model.EventTypeUpdatemember
	case *model.UpdateRole:
		return model.EventTypeUpdaterole
	case *model.UpdateRoom:
		return model.EventTypeUpdateroom
	case *model.UpdateForm:
		return model.EventTypeUpdateform
	case *model.UpdateAllows:
		return model.EventTypeUpdateallows
	case *model.UpdateChat:
		return model.EventTypeUpdatechat
	case *model.NewRoom:
		return model.EventTypeNewroom
	case *model.CreateInvite:
		return model.EventTypeCreateinvite
	case *model.DeleteInvite:
		return model.EventTypeDeleteinvite

	}
	panic("no matches found")
}
