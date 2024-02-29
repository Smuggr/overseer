package models

import (
	"gorm.io/gorm"
)


type User struct {
    gorm.Model
    Login           string `json:"login" gorm:"unique;not null"`
    Username        string `json:"username" gorm:"not null"`
    Password        string `json:"password" gorm:"not null"`
    PermissionLevel int    `json:"permission_level" gorm:"not null; default:1"`
}

func (User) TableName() string {
    return "users"
}


type Device struct {
    gorm.Model
    Login    string `json:"login" gorm:"unique;not null"`
    Username string `json:"username" gorm:"not null"`
    Password string `json:"password" gorm:"not null"`
}

func (Device) TableName() string {
    return "devices"
}