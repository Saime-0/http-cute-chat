package models

// INTERNAL models
// RECEIVED models
type CreateRole struct {
	RoleName      string `json:"role_name"`
	Color         string `json:"color"`
	Visible       bool   `json:"visible"`
	ManageRooms   bool   `json:"manage_rooms"`
	RoomID        int    `json:"room_id"`
	ManageChat    bool   `json:"manage_chat"`
	ManageRoles   bool   `json:"manage_roles"`
	ManageMembers bool   `json:"manage_members"`
}
type UpdateRole struct {
	RoleName      string `json:"role_name"`
	Color         string `json:"color"`
	Visible       bool   `json:"visible"`
	ManageRooms   bool   `json:"manage_rooms"`
	RoomID        int    `json:"room_id"`
	ManageChat    bool   `json:"manage_chat"`
	ManageRoles   bool   `json:"manage_roles"`
	ManageMembers bool   `json:"manage_members"`
}
type RoleInfo struct {
	ID       int    `json:"id"`
	RoleName string `json:"role_name"`
	Color    string `json:"color"`
}
type RoleData struct {
	ID            int    `json:"id"`
	RoleName      string `json:"role_name"`
	Color         string `json:"color"`
	Visible       bool   `json:"visible"`
	ManageRooms   bool   `json:"manage_rooms"`
	RoomID        int    `json:"room_id"`
	ManageChat    bool   `json:"manage_chat"`
	ManageRoles   bool   `json:"manage_roles"`
	ManageMembers bool   `json:"manage_members"`
}
type ListRolesInfo struct {
	Roles []RoleInfo `json:"roles"`
}
type ListRolesData struct {
	Roles []RoleData `json:"roles"`
}
type RoleID struct {
	ID int `json:"id"`
}
