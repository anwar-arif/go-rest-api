package model

type User struct {
	UserName  string   `json:"user_name" bson:"user_name"`
	Email     string   `json:"email" bson:"email"`
	Password  string   `json:"password" bson:"password"`
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

type GetByEmailRequest struct {
	Email string `json:"email"`
}

type UserResponse struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}
