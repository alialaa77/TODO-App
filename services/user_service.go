package services

import (
	"testing"
)

func TestParseDueDate(t *testing.T) {
	s := NewTodoService(nil)

	in := "2025-07-20T12:00:00Z"
	out, err := s.parseDueDate(&in)
	if err != nil || out == nil {
		t.Fatalf("expected valid due date, got err=%v", err)
	}

	past := "2025-07-10T12:00:00Z"
	_, err = s.parseDueDate(&past)
	if err == nil {
		t.Fatalf("expected error for past due date")
	}

	_, err = s.parseDueDate(nil)
	if err != nil {
		t.Fatalf("expected no error for nil due date")
	}

	invalid := "2025-07-20"
	_, err = s.parseDueDate(&invalid)
	if err == nil {
		t.Fatalf("expected parse error for invalid format")
	}
}

func TestValidatePriority(t *testing.T) {
	s := NewTodoService(nil)
	if !s.validatePriority("Low") || !s.validatePriority("medium") || !s.validatePriority("HIGH") {
		t.Fatalf("expected valid priorities")
	}
	if s.validatePriority("invalid") {
		t.Fatalf("expected invalid priority")
	}
}
