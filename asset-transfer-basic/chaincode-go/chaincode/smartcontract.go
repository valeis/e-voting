package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Voter struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HasVoted bool   `json:"hasVoted"`
}

type Candidate struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ElectionID string `json:"electionID"`
	Votes      int    `json:"votes"`
}

type Vote struct {
	VoterID     string `json:"VoterID"`
	CandidateID string `json:"candidateID"`
}

func (s *SmartContract) RegisterVoter(ctx contractapi.TransactionContextInterface, voterID string, name string) error {
	voter := Voter{
		ID:       voterID,
		Name:     name,
		HasVoted: false,
	}

	voterJSON, err := json.Marshal(voter)
	if err != nil {
		return fmt.Errorf("failed to marshal voter: %v", err)
	}

	return ctx.GetStub().PutState(voterID, voterJSON)
}

func (s *SmartContract) RegisterCandidate(ctx contractapi.TransactionContextInterface, candidateID string, name string, electionID string) error {
	candidate := Candidate{
		ID:         candidateID,
		Name:       name,
		ElectionID: electionID,
		Votes:      0,
	}

	candidateJSON, err := json.Marshal(candidate)
	if err != nil {
		return fmt.Errorf("failed to marshal candidate: %v", err)
	}

	return ctx.GetStub().PutState(candidateID, candidateJSON)
}

func (s *SmartContract) CastVote(ctx contractapi.TransactionContextInterface, voterID string, candidateID string) error {
	voterJSON, err := ctx.GetStub().GetState(voterID)
	if err != nil {
		return fmt.Errorf("failed to read voter: %v", err)
	}
	if voterJSON == nil {
		return fmt.Errorf("voter %s does not exist", voterID)
	}

	var voter Voter
	err = json.Unmarshal(voterJSON, &voter)
	if err != nil {
		return fmt.Errorf("failed to unmarshal voter: %v", err)
	}

	if voter.HasVoted {
		return fmt.Errorf("voter %s has already voted", voterID)
	}

	candidateJSON, err := ctx.GetStub().GetState(candidateID)
	if err != nil {
		return fmt.Errorf("failed to read candidate: %v", err)
	}
	if candidateJSON == nil {
		return fmt.Errorf("candidate %s does not exist", candidateID)
	}

	var candidate Candidate
	err = json.Unmarshal(candidateJSON, &candidate)
	if err != nil {
		return fmt.Errorf("failed to unmarshal candidate: %v", err)
	}

	candidate.Votes++
	candidateJSON, err = json.Marshal(candidate)
	if err != nil {
		return fmt.Errorf("failed to marshal candidate: %v", err)
	}
	err = ctx.GetStub().PutState(candidateID, candidateJSON)
	if err != nil {
		return fmt.Errorf("failed to update candidate: %v", err)
	}

	voter.HasVoted = true
	voterJSON, err = json.Marshal(voter)
	if err != nil {
		return fmt.Errorf("failed to marshal voter: %v", err)
	}
	err = ctx.GetStub().PutState(voterID, voterJSON)
	if err != nil {
		return fmt.Errorf("failed to update voter: %v", err)
	}

	return nil
}

func (s *SmartContract) GetVoteCount(ctx contractapi.TransactionContextInterface, candidateID string) (int, error) {
	candidateJSON, err := ctx.GetStub().GetState(candidateID)
	if err != nil {
		return 0, fmt.Errorf("failed to read candidate: %v", err)
	}
	if candidateJSON == nil {
		return 0, fmt.Errorf("candidate %s does not exist", candidateID)
	}

	var candidate Candidate
	err = json.Unmarshal(candidateJSON, &candidate)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal candidate: %v", err)
	}

	return candidate.Votes, nil
}

func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Candidate, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var candidates []*Candidate
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var candidate Candidate
		err = json.Unmarshal(queryResponse.Value, &candidate)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, &candidate)
	}

	return candidates, nil
}

func (s *SmartContract) GetCandidatesByElection(ctx contractapi.TransactionContextInterface, electionID string) ([]*Candidate, error) {
	if len(electionID) == 0 {
		return nil, fmt.Errorf("electionID cannot be empty")
	}

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	defer resultsIterator.Close()

	var candidates []*Candidate
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate through results: %v", err)
		}

		var candidate Candidate
		err = json.Unmarshal(queryResponse.Value, &candidate)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal candidate: %v", err)
		}

		if candidate.ElectionID == electionID {
			candidates = append(candidates, &candidate)
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates found for election ID: %s", electionID)
	}

	return candidates, nil
}
