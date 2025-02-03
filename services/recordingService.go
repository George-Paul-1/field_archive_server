package services

import (
	"context"
	"field_archive/server/entities"
	"field_archive/server/repositories"
	"fmt"
)

type RecordingService interface {
	GetByID(id int, ctx context.Context) (entities.Recording, error)
	ListItems(limit int, ctx context.Context) ([]entities.Recording, error)
	GetCount(ctx context.Context) (int, error)
}

type recordingService struct {
	repo repositories.RecordingRepository
}

func NewRecordingService(repo repositories.RecordingRepository) *recordingService {
	return &recordingService{repo: repo}
}

func (s *recordingService) GetByID(id int, ctx context.Context) (entities.Recording, error) {
	if id < 1 {
		return entities.Recording{}, fmt.Errorf("id must be no less than 1")
	}
	recording, err := s.repo.GetRowByID(id, ctx)
	if err != nil {
		return entities.Recording{}, fmt.Errorf("service: problem retrieving recording by ID, %w", err)
	}
	return recording, nil
}

func (s *recordingService) ListItems(limit int, ctx context.Context) ([]entities.Recording, error) {
	if limit < 1 {
		return []entities.Recording{}, fmt.Errorf("limit can't be less than 1")
	}
	recordings, err := s.repo.List(ctx, limit)
	if err != nil {
		return []entities.Recording{}, fmt.Errorf("service: problem retrieving list, %w", err)
	}
	return recordings, nil
}

func (s *recordingService) GetCount(ctx context.Context) (int, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("service: problem retrieving count, %w", err)
	}
	return count, nil

}
