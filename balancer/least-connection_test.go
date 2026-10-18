package balancer

import (
	"testing"
)

// Test 1: Verify that initServers correctly loads and parses the backend pool
func TestInitServers(t *testing.T) {
	pool, err := initServers()
	if err != nil {
		t.Fatalf("expected no error initializing servers, got %v", err)
	}

	expectedLen := 3
	if len(pool) != expectedLen {
		t.Errorf("expected %d servers in pool, got %d", expectedLen, len(pool))
	}

	// Verify the first server URL parsed correctly
	if pool[0].url.Host != "localhost:8000" {
		t.Errorf("expected host localhost:8000, got %s", pool[0].url.Host)
	}
}

// Test 2: Verify that NextServer accurately picks the server with the lowest count
func TestNextServerLeastConnections(t *testing.T) {
	// Construct a controlled pool with explicit connection counts
	pool := ServerPool{
		&Server{count: 5},
		&Server{count: 2}, // This should be picked because 2 is the minimum
		&Server{count: 10},
	}

	selected := pool.NextServer()
	if selected == nil {
		t.Fatalf("expected a server to be returned, got nil")
	}

	if selected.count != 2 {
		t.Errorf("expected server with count 2, got selected server with count %d", selected.count)
	}
}
