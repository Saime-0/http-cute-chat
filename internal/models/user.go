package models

// general user model
type User struct {
	ID          int    `json:"id"`
	Domain      string `json:"domain"`
	Name        string `json:"name"`
	AppSettings string `json:"app_settings"`
}

// RegisteredAt time.Time `json:"registeredAt"`
// LastVisitAt time.Time `json:"lastVisitAt"`
// Email string `json:"email"`

//  INTERNAL models
type UserDomain struct {
	Domain string `json:"domain"`
}

// RECEIVED models
type UpdateUserInput struct {
	Domain      string `json:"domain"`
	Name        string `json:"name"`
	AppSettings string `json:"app_settings"`
}

type CreateUserInput struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

// RESPONSE models
type UserInfo struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
}
