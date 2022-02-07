package model

type User struct {
	Name          string `json:"name" bson:"name"`
	LocalAddress  string `json:"local_address" bson:"local_address"`
	GlobalAddress string `json:"global_address" bson:"global_address"`
	Chatting      bool   `json:"chatting" bson:"chatting"`
}
type Users []User

type UserResponse struct {
	Users InnerUsers `json:"users"`
}

type InnerUser struct {
	Name     string `json:"name"`
	Chatting bool   `json:"chatting"`
}
type InnerUsers []InnerUser
