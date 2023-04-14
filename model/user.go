package model

type User struct {
	UserName  string   `json:"user_name" bson:"user_name"`
	Email     string   `json:"email" bson:"email"`
	Password  string   `json:"password" bson:"password"`
	Salt      string   `json:"salt" bson:"salt"`
	FirstName string   `json:"first_name,omitempty" bson:"first_name,omitempty"`
	Groups    []string `json:"groups,omitempty" bson:"groups,omitempty"`
	Contacts  []string `json:"contacts,omitempty" bson:"contacts,omitempty"`
	LastName  string   `json:"last_name,omitempty" bson:"last_name,omitempty"`
	IsActive  bool     `json:"is_active,omitempty" bson:"is_active,omitempty"`
	IsAdmin   bool     `json:"is_admin,omitempty" bson:"is_admin,omitempty"`
	UserID    string   `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Token     string   `json:"token,omitempty" bson:"token,omitempty"`
	Role      string   `json:"role,omitempty" bson:"role,omitempty"`
}

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

type GetUserByEmailResponse struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type AuthUserData struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active,omitempty" bson:"is_active,omitempty"`
	IsAdmin  bool   `json:"is_admin,omitempty" bson:"is_admin,omitempty"`
	UserID   string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Token    string `json:"token,omitempty" bson:"token,omitempty"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
}

type AuthUserResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    *AuthUserData `json:"data"`
}
