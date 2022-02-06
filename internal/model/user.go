package model

type User struct {
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}
type Users []User

type UserResponse struct {
	Users Users `json:"users"`
}
