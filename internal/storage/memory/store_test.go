package memory

import "testing"

func TestStoreAppend(t *testing.T) {
	store := New()
	store.SetCurrentTerm(3)
	store.SetVotedFor("node-2")
	store.Append(LogEntry{Index: 1, Term: 3, Data: []byte("hello")})

	if store.CurrentTerm() != 3 {
		t.Fatalf("unexpected term: %d", store.CurrentTerm())
	}
	if store.VotedFor() != "node-2" {
		t.Fatalf("unexpected vote target: %s", store.VotedFor())
	}
	if len(store.Entries()) != 1 {
		t.Fatalf("unexpected log size: %d", len(store.Entries()))
	}
}
