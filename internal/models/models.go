package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Login     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"not null;default:'viewer'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Device struct {
	ID            uint   `gorm:"primaryKey"`
	Type          string `gorm:"not null"`
	Vendor        string
	Model         string
	Serial        string `gorm:"unique;not null"`
	Location      string
	Status        string `gorm:"default:'active'"`
	NetworkNodeID *uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type NetworkNode struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	ParentID    *uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
