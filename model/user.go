package model

type User struct {
	UserID    string   `json:"user_id" bson:"user_id"`
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
	Token     string   `json:"token,omitempty" bson:"token"`
	Role      string   `json:"role,omitempty" bson:"role"`
}

type SignUpRequest struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

type GetUserByEmailResponse struct {
	UserName  string   `json:"user_name"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	Groups    []string `json:"groups"`
	Contacts  []string `json:"contacts"`
	LastName  string   `json:"last_name"`
	IsActive  bool     `json:"is_active"`
	IsAdmin   bool     `json:"is_admin"`
	Role      string   `json:"role"`
}

func (u *User) ToGetUserByEmailResponse() *GetUserByEmailResponse {
	return &GetUserByEmailResponse{
		UserName:  u.UserName,
		Email:     u.Email,
		FirstName: u.FirstName,
		Groups:    u.Groups,
		Contacts:  u.Contacts,
		LastName:  u.LastName,
		IsActive:  u.IsActive,
		IsAdmin:   u.IsAdmin,
		Role:      u.Role,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type LogOutRequest struct {
	Email string `json:"email"`
}

type AuthUserData struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	IsActive bool   `json:"is_active"`
	IsAdmin  bool   `json:"is_admin"`
	Token    string `json:"token" `
	Role     string `json:"role"`
}

func (au *AuthUserData) ToUser() *User {
	return &User{
		UserID:   au.UserID,
		UserName: au.UserName,
		Email:    au.Email,
		Role:     au.Role,
	}
}
