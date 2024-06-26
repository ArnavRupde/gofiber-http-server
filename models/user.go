package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model     `json:"-"` // `-` to hide fields generated by gorm in JSON
	Id             int        `json:"id" gorm:"primaryKey"`
	Name           string     `json:"name"`
	internal_id    string     // Small case for unexported field
	Username       string     `json:"-"` // `-` to hide fields generated by gorm in JSON
	HashedPassword string     `json:"-"` // `-` to hide fields generated by gorm in JSON
}
