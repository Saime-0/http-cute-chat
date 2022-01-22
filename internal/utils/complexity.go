package utils

import "github.com/saime-0/http-cute-chat/graph/generated"

func add(v int) func(childComplexity int) int {
	return func(childComplexity int) int {
		return childComplexity + v
	}
}
func mul(v int) func(childComplexity int) int {
	return func(childComplexity int) int {
		return childComplexity * v
	}
}

func MatchComplexity() *generated.ComplexityRoot {
	c := &generated.ComplexityRoot{}
	c.Chat.Members = add(3)
	c.Chat.Rooms = add(3)
	c.Room.Chat = add(2)
	c.Room.Allows = add(2)
	//c.Room.Messages = add(4)
	c.Messages.Messages = add(4)
	c.Members.Members = add(3)
	c.Chats.Chats = add(2)
	c.Units.Units = add(2)
	c.Users.Users = add(2)
	return c
}
