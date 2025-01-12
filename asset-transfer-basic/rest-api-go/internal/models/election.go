package models

import "gorm.io/gorm"

type Election struct {
	ID                   uint   `gorm:"primaryKey; autoIncrement:true;unique"`
	Title                string `gorm:"not null" json:"title"`
	Description          string `json:"description"`
	Type                 string `gorm:"not null" json:"type"`
	StartDate            string `gorm:"not null" json:"startDate"`
	EndDate              string `gorm:"not null" json:"endDate"`
	NumberOfAuthAttempts string `gorm:"not null" json:"numberOfAuthAttempts"`
	NumberOfCandidates   string `gorm:"not null" json:"numberOfCandidates"`
	NumberOfSelection    string `gorm:"not null" json:"numberOfSelection"`
	AuthMethod           string `gorm:"not null" json:"authMethod"`
}

type Candidate struct {
	gorm.Model
	FirstName  string `json:"firstName" gorm:"not null"`
	LastName   string `json:"lastName" gorm:"not null"`
	Age        int    `json:"age" gorm:"not null"`
	Party      string `json:"party" gorm:"not null"`
	Photo      string `json:"photo" gorm:"type:text"`
	ElectionID uint   `json:"electionID"`
}
