package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"rest-api-go/internal/models"
)

type UserRepository struct {
	dbClient any
}

func NewUserRepository(dbClient any) *UserRepository {
	return &UserRepository{dbClient: dbClient}
}

func (repo *UserRepository) GetUser(idnp string) (*models.User, error) {
	db, ok := repo.dbClient.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("invalid dbClient type: expected *gorm.DB")
	}
	var user models.User
	err := db.First(&user, "idnp = ?", idnp).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) IsUserRegistered(idnp string) (bool, error) {
	user, err := repo.GetUser(idnp)
	if err != nil {
		return false, err
	}
	return user.Registered, nil
}

func (repo *UserRepository) MarkUserAsRegistered(idnp string) error {
	db, ok := repo.dbClient.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid dbClient type: expected *gorm.DB")
	}
	existingUser := &models.User{}
	err := db.Where("idnp = ?", idnp).First(&existingUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with %s idnp not found", idnp)
		}
		return err
	}
	existingUser.Registered = true
	err = db.Save(&existingUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) AddUser(user *models.User) error {
	db, ok := repo.dbClient.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid dbClient type: expected *gorm.DB")
	}
	return db.Create(user).Error
}
