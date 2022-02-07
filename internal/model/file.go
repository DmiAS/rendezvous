package model

type FileMeta struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}
type FilesMeta []FileMeta

type FileRecord struct {
	User  string    `json:"user" bson:"user"`
	Files FilesMeta `json:"files" bson:"files"`
}
