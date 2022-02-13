package model

type User struct {
	Name          string `json:"name" bson:"name"`
	LocalAddress  string `json:"local_address" bson:"local_address"`
	GlobalAddress string `json:"global_address" bson:"global_address"`
	Blocked       bool   `json:"chatting" bson:"blocked"`
}
type Users []User

type UserResponse struct {
	Users InnerUsers `json:"users"`
}

type InnerUsers struct {
	Names []string `json:"names"`
}
