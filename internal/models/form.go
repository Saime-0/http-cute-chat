package models

type FormField struct {
	Key      string `json:"key"`      // unique, but used for grouping radiobutton type fields
	Type     string `json:"type"`     // email, data(past|future), link(by domain(ex:https://youtube.com/)),
	Optional bool   `json:"optional"` // omitempty
	Length   int    `json:"length"`   // omitempty
	Value    string `json:"value"`
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
