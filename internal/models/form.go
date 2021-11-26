package models

import "github.com/saime-0/http-cute-chat/internal/api/rules"

type FormField struct {
	Key      string          `json:"key"`      // unique, but used for grouping radiobutton type fields
	Type     rules.FieldType `json:"type"`     // email, date(past|future), link(by domain(ex:https://youtube.com/)),
	Optional bool            `json:"optional"` // omitempty
	Length   int             `json:"length"`   // omitempty
	Items    []string        `json:"items"`    // field for radiobutton and similar
}

type FormPattern struct {
	Fields []FormField `json:"fields"`
}

type FormChoice struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type FormCompleted struct {
	Input []FormChoice `json:"input"`
}
