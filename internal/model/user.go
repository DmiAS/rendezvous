package model

type User struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
type Users []User
