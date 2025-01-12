package controller

import (
	"errors"
	"fmt"
	"rest-api-go/internal/models"
)

type IUserRepository interface {
	GetUser(idnp string) (*models.User, error)
	IsUserRegistered(idnp string) (bool, error)
	MarkUserAsRegistered(idnp string) error
	AddUser(user *models.User) error
}

type ChainCodeProxy struct {
	userRepo IUserRepository
}

// interface to build a User
type UserBuilder interface {
	SetIDNP(string) UserBuilder
	SetName(string) UserBuilder
	Validate() error
	Build() (*models.User, error)
}

// concrete builder for building user
type VoterUserBuilder struct {
	user *models.User
	repo IUserRepository
}

// creates a new instance of VoterUserBuilder
func NewVoterUserBuilder(userRepo IUserRepository) *VoterUserBuilder {
	return &VoterUserBuilder{
		user: &models.User{
			Registered: true,
			Role:       "user",
		},
		repo: userRepo,
	}
}

func (b *VoterUserBuilder) SetIDNP(idnp string) UserBuilder {
	b.user.Idnp = idnp
	return b
}

func (b *VoterUserBuilder) SetName(name string) UserBuilder {
	b.user.Name = name
	return b
}

func (b *VoterUserBuilder) Validate() error {
	if b.user.Idnp == "" {
		return errors.New("IDNP cannot be empty")
	}

	isRegistered, _ := b.repo.IsUserRegistered(b.user.Idnp)
	//if err != nil {
	//	return fmt.Errorf("error checking user registration: %v", err)
	//}
	if isRegistered {
		return fmt.Errorf("user with IDNP %s already registered", b.user.Idnp)
	}
	return nil
}

func (b *VoterUserBuilder) Build() (*models.User, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}

	if err := b.repo.AddUser(b.user); err != nil {
		return nil, fmt.Errorf("can't register new user in local database: %v", err)
	}

	return b.user, nil
}

func NewChaincodeProxy(userRepo IUserRepository) *ChainCodeProxy {
	return &ChainCodeProxy{userRepo: userRepo}
}

func (proxy *ChainCodeProxy) ValidateAndForward(function string, args []string, forward func() (string, error)) (string, error) {
	if function == "RegisterVoter" {
		if len(args) == 0 {
			return "", errors.New("missing user IDNP")
		}
		builder := NewVoterUserBuilder(proxy.userRepo)
		_, err := builder.
			SetIDNP(args[0]).
			SetName(args[1]).
			Build()

		if err != nil {
			return "", err
		}

		return "voter registered successfully", nil

		//user := models.User{}
		//user.Registered = true
		//user.Idnp = args[0]
		//user.Name = args[1]
		//user.Role = "user"
		//
		//userIDNP := args[0]

		//isUserRegistered, _ := proxy.userRepo.IsUserRegistered(userIDNP)
		//if isUserRegistered {
		//	return "", fmt.Errorf("user with %s idnp already registered", userIDNP)
		//}
		//err := proxy.userRepo.AddUser(&user)
		//if err != nil {
		//	fmt.Println("can't register new user in local database")
		//}
	}
	return forward()
}
