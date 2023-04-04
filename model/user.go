package model

type User struct {
	UserName string `json:"userName" bson:"userName"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}
