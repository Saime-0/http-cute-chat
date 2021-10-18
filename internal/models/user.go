package models

// general user model
type User struct {
	ID          int              `json:"id"`
	Domain      string           `json:"domain"`
	Name        string           `json:"name"`
	AppSettings string           `json:"app_settings"`
	Sessions    []RefreshSession `json:"sessions"`
}
type RefreshSession struct {
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	Exp          int    `json:"exp"`
	CreatedAt    int    `json:"created_at"`
}

// RegisteredAt time.Time `json:"registeredAt"`
// LastVisitAt time.Time `json:"lastVisitAt"`
// Email string `json:"email"`

// INTERNAL models
// RECEIVED models
// todo: input data to refresh_sessions
type CreateUser struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
}
type UserInput struct { // pass & email
	Domain string `json:"domain"`
	Name   string `json:"name"`
}
type UpdateUserData struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
}
type UpdateUserSettings struct {
	AppSettings string `json:"app_settings"`
}

type TokenForRefreshPair struct {
	RefreshToken string `json:"refresh_token"`
}
type UserName struct {
	Name string `json:"name"`
}

// RESPONSE models
type UserInfo struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
}
type UserData struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
}
type UserSettings struct {
	AppSettings string `json:"app_settings"`
}
type FreshTokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type ListUsers struct {
	Users []UserInfo `json:"users"`
}

// ? todo: search user by name
// ? todo: get chats controlled by user
// todo: user sign im model
