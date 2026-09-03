package memory

import (
	"database/sql"
	"fmt"
	"time"
)

type MemoryScope string

const (
	ScopeGlobal  MemoryScope = "GLOBAL"
	ScopeProject MemoryScope = "PROJECT"
	ScopeSession MemoryScope = "SESSION"
)

type MemoryEntry struct {
	ID        string
	Scope     MemoryScope
	Tags      []string
	Content   string
	CreatedAt time.Time
}

type Store struct {
	db *sql.DB
}

// NewStore initializes the SQLite memory index (stubbed interface for now)
func NewStore(dbPath string) (*Store, error) {
	// Future: db, err := sql.Open("sqlite", dbPath)
	return &Store{}, nil
}

// StoreEntry saves an entry via structured SQLite insertion
func (s *Store) StoreEntry(entry *MemoryEntry) error {
	fmt.Printf("[Memory] Storing indexed entry [%s]: %v\n", entry.Scope, entry.Tags)
	// SQL INSERT logic here
	return nil
}

// RetrieveByTags fetches memory entries using indexed keyword/tag matching (no vectors yet)
func (s *Store) RetrieveByTags(scope MemoryScope, tags []string) ([]MemoryEntry, error) {
	fmt.Printf("[Memory] Retrieving %s memory for tags: %v\n", scope, tags)
	// SQL SELECT logic using index/tags
	return []MemoryEntry{}, nil
}
