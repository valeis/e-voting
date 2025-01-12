package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"rest-api-go/internal/models"
)

type IElectionRepository interface {
	RegisterElection(election *models.Election) error
	GetAllElection() ([]models.Election, error)
	GetElectionById(id uint) (*models.Election, error)
}

type ElectionRepository struct {
	dbClient any
}

func NewElectionRepository(dbClient any) *ElectionRepository {
	return &ElectionRepository{dbClient: dbClient}
}

func (repo *ElectionRepository) RegisterElection(election *models.Election) error {
	db, ok := repo.dbClient.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid dbClient type: expected *gorm.DB")
	}
	res := db.Create(&election)
	if res.Error != nil {
		return nil
	}
	return nil
}

func (repo *ElectionRepository) GetAllElection() ([]models.Election, error) {
	db, ok := repo.dbClient.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("invalid dbClient type: expected *gorm.DB")
	}

	var elections []models.Election
	res := db.Find(&elections)
	if res.Error != nil {
		return nil, res.Error
	}
	return elections, nil
}

func (repo *ElectionRepository) GetElectionById(id uint) (*models.Election, error) {
	db, ok := repo.dbClient.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("invalid dbClient type: expected *gorm.DB")
	}

	var election models.Election
	result := db.First(&election, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("election with id %d not found", id)
		}
		return nil, result.Error
	}
	return &election, nil
}
