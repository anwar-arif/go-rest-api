package model

type User struct {
	UserName  string   `json:"user_name" bson:"user_name"`
	Email     string   `json:"email" bson:"email"`
	Password  string   `json:"password" bson:"password"`
	Salt      string   `json:"salt" bson:"salt"`
	FirstName string   `json:"first_name,omitempty" bson:"first_name"`
	Groups    []string `json:"groups,omitempty" bson:"groups"`
	Contacts  []string `json:"contacts,omitempty" bson:"contacts"`
	LastName  string   `json:"last_name,omitempty" bson:"last_name"`
	IsActive  bool     `json:"is_active,omitempty" bson:"is_active"`
	IsAdmin   bool     `json:"is_admin,omitempty" bson:"is_admin"`
	UserID    string   `json:"user_id,omitempty" bson:"user_id"`
	Token     string   `json:"token,omitempty" bson:"token"`
	Role      string   `json:"role,omitempty" bson:"role"`
}

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

type GetUserByEmailResponse struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type AuthUserData struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	IsActive bool   `json:"is_active"`
	IsAdmin  bool   `json:"is_admin"`
	UserID   string `json:"user_id" `
	Token    string `json:"token" `
	Role     string `json:"role"`
}

func (au *AuthUserData) ToUser() *User {
	return &User{
		UserName: au.UserName,
		Email:    au.Email,
		Role:     au.Role,
	}
}
