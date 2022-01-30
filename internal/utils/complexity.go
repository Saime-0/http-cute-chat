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
	c.Chat.Members = add(4)
	c.Chat.Rooms = add(4)
	c.Room.Chat = add(3)
	c.Room.Allows = add(3)
	//c.Room.Messages = add(5)
	c.Messages.Messages = add(5)
	c.Members.Members = add(4)
	c.Chats.Chats = add(3)
	c.Units.Units = add(3)
	c.Users.Users = add(3)
	return c
}
