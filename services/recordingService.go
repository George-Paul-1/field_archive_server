package services

import (
	"context"
	"field_archive/server/entities"
	"field_archive/server/repositories"
	"fmt"
)

type RecordingService struct {
	repo repositories.RecordingRepository
}

func NewRecordingService(repo repositories.RecordingRepository) *RecordingService {
	return &RecordingService{repo: repo}
}

func (s *RecordingService) GetByID(id int, ctx context.Context) (entities.Recording, error) {
	if id < 1 {
		return entities.Recording{}, fmt.Errorf("id must be no less than 1")
	}
	recording, err := s.repo.GetRowByID(id, ctx)
	if err != nil {
		return entities.Recording{}, fmt.Errorf("service: problem retrieving recording by ID, %w", err)
	}
	return recording, nil
}
