package model

type UserTable struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	LastLogin string `json:"last_login"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

func (UserTable) TableName() string {
	return "users"
}

type UserUpdatePassword struct {
	ID          int64  `json:"id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
