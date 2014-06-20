package models

import (
	"testing"
)

func TestHasVoted(t *testing.T) {

	u := &User{}
	if u.HasVoted(1) {
		t.Error("The user has not voted yet")
	}

	u.AddVote(1)
	if !u.HasVoted(1) {
		t.Error("The user should have voted")
	}
}
