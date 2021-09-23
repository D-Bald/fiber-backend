package model

// Role identified by `Tag` field and represented by `Name` field
type Role struct {
	Tag  string `bson:"tag,omitempty" json:"tag" xml:"tag" form:"tag"`
	Name string `bson:"name,omitempty" json:"name" xml:"name" form:"name"`
}

// Only `Name` field that can be updated through API endpoint
type RoleUpdate struct {
	Name string `bson:"name,omitempty" json:"name" xml:"name" form:"name"`
}
