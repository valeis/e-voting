package service

import (
	"fmt"
	"rest-api-go/internal/models"
	"rest-api-go/internal/repository"
)

type ElectionService interface {
	RegisterElection(election *models.Election) error
	GetAllElections() ([]models.Election, error)
	GetElectionById(id uint) (*models.Election, error)
	RegisterCandidates(electionID uint, candidates []models.Candidate) error
	GetCandidates(electionId uint) ([]models.Candidate, error)
}

type ElectionServiceImpl struct {
	ElectionRepository   repository.IElectionRepository
	CandidatesRepository repository.ICandidateRepository
}

func NewElectionServiceImpl(electionRepository repository.IElectionRepository, candidateRepository repository.ICandidateRepository) (service ElectionService) {
	return &ElectionServiceImpl{ElectionRepository: electionRepository, CandidatesRepository: candidateRepository}
}

func (electionRepo *ElectionServiceImpl) RegisterElection(election *models.Election) (err error) {
	err = electionRepo.ElectionRepository.RegisterElection(election)
	if err != nil {
		return err
	}
	return nil
}

func (electionRepo *ElectionServiceImpl) GetAllElections() ([]models.Election, error) {
	elections, err := electionRepo.ElectionRepository.GetAllElection()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch elections: %v", err)
	}
	return elections, nil
}

func (electionRepo *ElectionServiceImpl) GetElectionById(id uint) (*models.Election, error) {
	election, err := electionRepo.ElectionRepository.GetElectionById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch election: %v", err)
	}
	return election, nil
}

func (electionRepo *ElectionServiceImpl) RegisterCandidates(electionID uint, candidates []models.Candidate) error {
	election, err := electionRepo.ElectionRepository.GetElectionById(electionID)
	if err != nil {
		return fmt.Errorf("could not find election with id %d: %w", electionID, err)
	}

	for i := range candidates {
		candidates[i].ElectionID = election.ID
	}

	if err := electionRepo.CandidatesRepository.RegisterCandidates(candidates); err != nil {
		return fmt.Errorf("failed to register candidates: %w", err)
	}
	return nil
}

func (electionRepo *ElectionServiceImpl) GetCandidates(electionId uint) ([]models.Candidate, error) {
	return electionRepo.CandidatesRepository.GetCandidatesByElectionId(electionId)
}
