package repository

import (
	"fmt"
	"gorm.io/gorm"
	"rest-api-go/internal/models"
)

type ICandidateRepository interface {
	RegisterCandidates(candidates []models.Candidate) error
	GetCandidatesByElectionId(electionId uint) ([]models.Candidate, error)
}

type CandidateRepository struct {
	dbClient any
}

func NewCandidateRepository(dbClient any) *CandidateRepository {
	return &CandidateRepository{dbClient: dbClient}
}

func (repo *CandidateRepository) RegisterCandidates(candidates []models.Candidate) error {
	db, ok := repo.dbClient.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid dbClient type: expected *gorm.DB")
	}
	if err := db.Create(&candidates).Error; err != nil {
		return fmt.Errorf("failed to register candidates: %w", err)
	}
	return nil
}

func (repo *CandidateRepository) GetCandidatesByElectionId(electionId uint) ([]models.Candidate, error) {
	db, ok := repo.dbClient.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("invalid dbClient type: expected *gorm.DB")
	}

	var candidates []models.Candidate
	res := db.Where("election_id = ?", electionId).Find(&candidates)
	if res.Error != nil {
		return nil, res.Error
	}

	return candidates, nil
}
