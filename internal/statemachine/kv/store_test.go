package kv

import "testing"

func TestStorePutAndGet(t *testing.T) {
	store := New()
	store.Put("demo", "value")

	got, ok := store.Get("demo")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if got != "value" {
		t.Fatalf("unexpected value: %s", got)
	}
}
