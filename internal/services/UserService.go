package svc

import "github.com/Desgue/ttracker-api/internal/domain"

// User service that handles business logic before inserting user into the database

type UserService struct {
	store domain.UserStorage
}

func NewUserService(store domain.UserStorage) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) CreateUser(cognitoId string) error {
	if err := s.store.CreateUser(cognitoId); err != nil {
		return err
	}
	return nil
}
