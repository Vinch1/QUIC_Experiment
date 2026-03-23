package memory

import "sync"

type LogEntry struct {
	Index uint64
	Term  uint64
	Data  []byte
}

type Store struct {
	mu          sync.RWMutex
	currentTerm uint64
	votedFor    string
	logs        []LogEntry
}

func New() *Store {
	return &Store{
		logs: make([]LogEntry, 0),
	}
}

func (s *Store) CurrentTerm() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentTerm
}

func (s *Store) SetCurrentTerm(term uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentTerm = term
}

func (s *Store) VotedFor() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.votedFor
}

func (s *Store) SetVotedFor(nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.votedFor = nodeID
}

func (s *Store) Append(entry LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logs = append(s.logs, entry)
}

func (s *Store) Entries() []LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]LogEntry, len(s.logs))
	copy(result, s.logs)
	return result
}
