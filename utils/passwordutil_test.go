package passwordutil

import (
	"testing"
)

// TestGeneratePasswordHash should return valid has for the given input string
func TestValidGeneratePasswordHash(t *testing.T) {
	expectedHash := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	actualHash, err := GeneratePasswordHash("angryMonkey")
	t.Log(actualHash)
	if expectedHash != actualHash {
		t.Errorf("Expected %s but actual %s", expectedHash, actualHash)
	}
	if err != nil {
		t.Errorf("Failed to generate password hash for valid string")
	}
}

// TestInvalidPasswordHash should return error if blank password is used
func TestInvalidPasswordHash(t *testing.T) {
	actualHash, err := GeneratePasswordHash("  ")
	t.Log(actualHash)
	if err == nil {
		t.Errorf("Generated hash for empty string")
	}
}
