package service

import (
	"context"

	"github.com/rs/zerolog"
)

type Storage interface {
	GetDataFromDbByID(ctx context.Context, id uint64) (string, error)
}

type Service interface {
	GetDataFromServiceLayerByID(ctx context.Context, id uint64) (string, error)
}

type service struct {
	log     zerolog.Logger
	storage Storage
}

func (s *service) GetDataFromServiceLayerByID(ctx context.Context, id uint64) (string, error) {
	return s.storage.GetDataFromDbByID(ctx, id)
}

func New(log zerolog.Logger, storage Storage) Service {
	return &service{
		log:     log,
		storage: storage,
	}
}
