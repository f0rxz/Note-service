package storage

import (
	"errors"
	"note-service/models"
	"sync"
)

type MemoryStorage struct {
	notes map[string]*models.Note
	mu    sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		notes: make(map[string]*models.Note),
	}
}

func (s *MemoryStorage) GetAll() []*models.Note {
	s.mu.Lock()
	defer s.mu.Unlock()
	var notes []*models.Note
	for _, note := range s.notes {
		notes = append(notes, note)
	}
	return notes
}

func (s *MemoryStorage) Create(note *models.Note) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notes[note.ID] = note
}

func (s *MemoryStorage) Get(id string) (*models.Note, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	note, found := s.notes[id]
	return note, found
}

func (s *MemoryStorage) Update(id string, note *models.Note) (*models.Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, found := s.notes[id]; !found {
		return nil, errors.New("note not found")
	}
	s.notes[id] = note
	return note, nil
}

func (s *MemoryStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, found := s.notes[id]; !found {
		return errors.New("note not found")
	}
	delete(s.notes, id)
	return nil
}
