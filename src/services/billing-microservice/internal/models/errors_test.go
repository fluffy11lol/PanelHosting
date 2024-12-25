package models

import (
	"errors"
	"testing"
)

func TestErrorFailToInsert(t *testing.T) {
	expected := "failed to insert"

	if errors.Is(ErrorFailToInsert, errors.New(expected)) {
		t.Errorf("ErrorFailToInsert does not match expected value. Got: %v, Want: %v", ErrorFailToInsert, expected)
	}
}

func TestErrorFailToInsertMessage(t *testing.T) {
	expected := "failed to insert"

	if ErrorFailToInsert.Error() != expected {
		t.Errorf("Unexpected error message. Got: %s, Want: %s", ErrorFailToInsert.Error(), expected)
	}
}
