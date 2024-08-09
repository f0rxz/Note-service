package storage

import (
	"database/sql"
	"fmt"
	"log"
	"note-service/models"
	"sort"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db         *sql.DB
	ticker     *time.Ticker
	wg         sync.WaitGroup
	mu         sync.Mutex
	stop       chan bool
	keyNotes   []int64
	nextNoteId int64
	notes      map[int64]*models.Note
	writeCache map[int64][]*models.Note
}

func NewStorage(dataSourceName string, interval time.Duration) (*Storage, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	s := &Storage{
		db:         db,
		ticker:     time.NewTicker(interval),
		stop:       make(chan bool),
		keyNotes:   make([]int64, 0),
		nextNoteId: 1,
		notes:      make(map[int64]*models.Note),
		writeCache: make(map[int64][]*models.Note),
	}
	s.LoadNotes()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.ticker.C:
				if err := s.SaveNotes(); err != nil {
					log.Printf("Error saving notes: %v", err)
				}
			case <-s.stop:
				return
			}
		}
	}()

	return s, nil
}

func (s *Storage) LoadNotes() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query("SELECT id, title, content FROM notes")
	if err != nil {
		return fmt.Errorf("error querying notes: %w", err)
	}
	defer rows.Close()

	s.keyNotes = make([]int64, 0)
	s.notes = make(map[int64]*models.Note)

	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content); err != nil {
			return fmt.Errorf("error scanning note: %w", err)
		}

		s.notes[note.ID] = &note
		s.keyNotes = append(s.keyNotes, note.ID)
		s.nextNoteId = note.ID + 1
	}

	return nil
}

func (s *Storage) SaveNotes() error {
	s.mu.Lock()
	writeCache := s.writeCache
	s.writeCache = make(map[int64][]*models.Note)
	s.mu.Unlock()

	// Begin the transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	// Defer the rollback to ensure it's called if something fails
	defer func() {
		rollbackChanges := func() {
			tx.Rollback()
			s.mu.Lock()
			for id, noteChanges := range writeCache {
				s.writeCache[id] = append(noteChanges, s.writeCache[id]...)
			}
			s.mu.Unlock()
		}

		if p := recover(); p != nil {
			rollbackChanges() // Rollback in case of panic
			panic(p)          // Rethrow panic after rollback
		} else if err != nil {
			// Rollback in case of error
			rollbackChanges()
		} else {
			err = tx.Commit() // Commit if no error occurred
		}
	}()

	for id, noteChanges := range writeCache {
		for _, noteChange := range noteChanges {
			if noteChange == nil {
				_, err = tx.Exec("DELETE FROM notes WHERE id = ?", id)
				if err != nil {
					return fmt.Errorf("error deleting note with ID %s: %w", id, err)
				}
				log.Printf("Deleted note with ID %s", id)
			} else {
				var exists bool
				err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM notes WHERE id = ?)", id).Scan(&exists)
				if err != nil {
					return fmt.Errorf("error checking existence of note with ID %s: %w", id, err)
				}

				if exists {
					_, err = tx.Exec("UPDATE notes SET title = ?, content = ? WHERE id = ?", noteChange.Title, noteChange.Content, noteChange.ID)
					if err != nil {
						return fmt.Errorf("error updating note with ID %s: %w", id, err)
					}
					log.Printf("Updated note with ID %s", id)
				} else {
					_, err = tx.Exec("INSERT INTO notes (id, title, content) VALUES (?, ?, ?)", noteChange.ID, noteChange.Title, noteChange.Content)
					if err != nil {
						return fmt.Errorf("error inserting note with ID %s: %w", id, err)
					}
					log.Printf("Inserted new note with ID %s", id)
				}
			}
		}
	}

	return err
}

func (s *Storage) GetNote(id int64) *models.Note {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.notes[id]
}

func (s *Storage) GetNotesRange(offset, limit int64) []*models.Note {
	s.mu.Lock()
	defer s.mu.Unlock()

	notes := make([]*models.Note, 0)
	keyNotesLength := int64(len(s.keyNotes))

	if offset < keyNotesLength {
		cuttingLength := min(keyNotesLength, offset+limit)

		for _, id := range s.keyNotes[offset:cuttingLength] {
			notes = append(notes, s.notes[id])
		}
	}

	return notes
}

func binarySearch(arr []int64, target int64) int {
	left, right := 0, len(arr)-1

	for left <= right {
		mid := left + (right-left)/2
		if arr[mid] == target {
			return mid
		}
		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

func (s *Storage) UpdateNote(id int64, note *models.Note) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id == 0 {
		id = s.nextNoteId
	}
	_, noteExists := s.notes[id]
	if note == nil {
		if noteExists {
			i := binarySearch(s.keyNotes, id)
			if i != -1 {
				newLen := len(s.keyNotes) - 1
				s.keyNotes[i] = s.keyNotes[newLen]
				s.keyNotes = s.keyNotes[:newLen]
				sort.Slice(s.keyNotes, func(i, j int) bool {
					return s.keyNotes[i] < s.keyNotes[j]
				})
			}
			delete(s.notes, id)
			s.writeCache[id] = append(s.writeCache[id], nil)
		}
	} else {
		note.ID = id
		if !noteExists {
			s.nextNoteId++
			s.keyNotes = append(s.keyNotes, id)
			sort.Slice(s.keyNotes, func(i, j int) bool {
				return s.keyNotes[i] < s.keyNotes[j]
			})
		}
		s.notes[id] = note
		s.writeCache[id] = append(s.writeCache[id], note)
	}
}

func (s *Storage) CreateNote(title, content string) {
	s.UpdateNote(0, &models.Note{
		Title:   title,
		Content: content,
	})
}

func (s *Storage) DeleteNote(id int64) {
	s.UpdateNote(id, nil)
}

func (s *Storage) EditNote(id int64, title, content string) {
	s.UpdateNote(id, &models.Note{
		Title:   title,
		Content: content,
	})
}

func (s *Storage) Close() error {
	s.ticker.Stop()
	close(s.stop)
	s.wg.Wait()
	return s.db.Close()
}
