package model

const TagTableName = "Tag"

type Tag struct {
	Id   int64  `db:"id" json:"id" struct:"id"`
	Name string `db:"name" json:"name" struct:"name"`
}

func NewTag(name string) *Tag {
	return &Tag{
		Name: name,
	}
}
