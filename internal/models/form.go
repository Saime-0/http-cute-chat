package models

type FormField struct {
	Key      string `json:"key"`      // unique, but used for grouping radiobutton type fields
	Type     string `json:"type"`     // email, date(past|future), link(by domain(ex:https://youtube.com/)),
	Optional bool   `json:"optional"` // omitempty
	Length   int    `json:"length"`   // omitempty
	Value    string `json:"value"`    // field for radiobutton and similar
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

var (
	CommonFormPattern = &FormPattern{
		Fields: []FormField{
			{
				Key:      "text",
				Type:     "text",
				Optional: false,
				Length:   0,
				Value:    "",
			},
		},
	}
)
