package models

type User struct {
	ID          int    `json:"id"`
	Domain      string `json:"domain"`
	Name        string `json:"name"`
	AppSettings string `json:"app_settings"`
}

type UserDomain struct {
	Domain string `json:"domain"`
}

type UpdateUserInput struct {
	Domain      string `json:"domain"`
	Name        string `json:"name"`
	AppSettings string `json:"app_settings"`
}
