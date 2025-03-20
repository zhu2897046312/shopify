package service

import (
	"shopify/repository"
)

type Service struct {
	repoFactory *repository.RepositoryFactory
}

func NewService(repoFactory *repository.RepositoryFactory) *Service {
	return &Service{
		repoFactory: repoFactory,
	}
} 