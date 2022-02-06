package model

type User struct {
	Name          string `json:"name" bson:"name"`
	LocalAddress  string `json:"local_address" bson:"local_address"`
	GlobalAddress string `json:"global_address" bson:"global_address"`
}
type Users []User

type UserResponse struct {
	Users Users `json:"users"`
}
