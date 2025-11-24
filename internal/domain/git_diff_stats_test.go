package domain

import "testing"

func TestGitDiffStats_Creation(t *testing.T) {
	// arrange
	linesAdded := 42
	linesRemoved := 17

	// act
	stats := GitDiffStats{
		LinesAdded:   linesAdded,
		LinesRemoved: linesRemoved,
	}

	// assert
	if stats.LinesAdded != linesAdded {
		t.Errorf("LinesAdded = %d, want %d", stats.LinesAdded, linesAdded)
	}
	if stats.LinesRemoved != linesRemoved {
		t.Errorf("LinesRemoved = %d, want %d", stats.LinesRemoved, linesRemoved)
	}
}

func TestGitDiffStats_ZeroValues(t *testing.T) {
	// act
	stats := GitDiffStats{
		LinesAdded:   0,
		LinesRemoved: 0,
	}

	// assert
	if stats.LinesAdded != 0 {
		t.Errorf("LinesAdded = %d, want 0", stats.LinesAdded)
	}
	if stats.LinesRemoved != 0 {
		t.Errorf("LinesRemoved = %d, want 0", stats.LinesRemoved)
	}
}
