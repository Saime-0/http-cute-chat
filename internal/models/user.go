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
	Exp          int64  `json:"exp"`
	CreatedAt    int64  `json:"created_at"`
}

type CustomClaims struct {
	UserID    int   `json:"user-id"`
	ExpiresAt int64 `json:"exp"`
}

// RegisteredAt time.Time `json:"registeredAt"`
// LastVisitAt time.Time `json:"lastVisitAt"`
// Email string `json:"email"`

// INTERNAL models
// RECEIVED models
// todo: input data to refresh_sessions
type CreateUser struct {
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserInput struct { // pass & email
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UpdateUserData struct {
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UpdateUserSettings struct {
	AppSettings string `json:"app_settings"`
}
type TokenForRefreshPair struct {
	RefreshToken string `json:"refresh_token"`
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
	Email  string `json:"email"`
}
type UserSettings struct {
	AppSettings string `json:"app_settings"`
}
type FreshTokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type ListUserInfo struct {
	Users []UserInfo `json:"users"`
}

// ? todo: search user by name
// ? todo: get chats controlled by user
// todo: user sign im model
