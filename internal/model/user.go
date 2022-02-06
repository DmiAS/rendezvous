package model

type User struct {
	Name          string `bson:"name"`
	LocalAddress  string `bson:"local_address"`
	GlobalAddress string `bson:"global_address"`
}
type Users []User

type UserResponse struct {
	Users Users `json:"users"`
}
