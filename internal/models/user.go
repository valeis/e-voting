package models

type User struct {
	ID         uint   `gorm:"primaryKey; autoIncrement:true;unique"`
	Idnp       string `gorm:"not null;unique" json:"idnp"`
	Name       string `gorm:"not null" json:"name"`
	Role       string `json:"role"`
	Registered bool   `json:"registered"`
}
