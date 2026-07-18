package api

import (
	"testing"
)

func TestNewStudentHandler(t *testing.T) {
	handler := NewStudentHandler(nil, nil)

	if handler == nil {
		t.Fatal("expected handler")
	}

	if handler.academicService != nil {
		t.Fatal("expected nil academic service")
	}

	if handler.ticketService != nil {
		t.Fatal("expected nil ticket service")
	}
}
