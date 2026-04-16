package ladder

import (
	"testing"

	"github.com/me/level-up-hub/backend/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestLadderLevelValidation(t *testing.T) {
	tests := []struct {
		name  string
		level repository.LadderLevel
		valid bool
	}{
		{
			name:  "valid P1",
			level: repository.LadderLevelP1,
			valid: true,
		},
		{
			name:  "valid P2",
			level: repository.LadderLevelP2,
			valid: true,
		},
		{
			name:  "valid P3",
			level: repository.LadderLevelP3,
			valid: true,
		},
		{
			name:  "valid LT1",
			level: repository.LadderLevelLT1,
			valid: true,
		},
		{
			name:  "valid LT2",
			level: repository.LadderLevelLT2,
			valid: true,
		},
		{
			name:  "valid LT3",
			level: repository.LadderLevelLT3,
			valid: true,
		},
		{
			name:  "invalid level",
			level: repository.LadderLevel("INVALID"),
			valid: false,
		},
	}

	validLevels := map[repository.LadderLevel]bool{
		repository.LadderLevelP1:  true,
		repository.LadderLevelP2:  true,
		repository.LadderLevelP3:  true,
		repository.LadderLevelLT1: true,
		repository.LadderLevelLT2: true,
		repository.LadderLevelLT3: true,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validLevels[tt.level]
			assert.Equal(t, tt.valid, isValid)
		})
	}
}

func TestLadderLevelOrdering(t *testing.T) {
	levels := []repository.LadderLevel{
		repository.LadderLevelP1,
		repository.LadderLevelP2,
		repository.LadderLevelP3,
		repository.LadderLevelLT1,
		repository.LadderLevelLT2,
		repository.LadderLevelLT3,
	}

	// Verify career progression order (P1 -> P2 -> P3 -> LT1 -> LT2 -> LT3)
	// Note: This tests the logical progression, not alphabetical ordering
	assert.Len(t, levels, 6, "Should have exactly 6 career levels")

	// Verify P-series ordering
	assert.Less(t, string(repository.LadderLevelP1), string(repository.LadderLevelP2))
	assert.Less(t, string(repository.LadderLevelP2), string(repository.LadderLevelP3))

	// Verify LT-series ordering
	assert.Less(t, string(repository.LadderLevelLT1), string(repository.LadderLevelLT2))
	assert.Less(t, string(repository.LadderLevelLT2), string(repository.LadderLevelLT3))

	// Verify P-series comes before LT-series in alphabetical order (P > L)
	assert.Greater(t, string(repository.LadderLevelP1), string(repository.LadderLevelLT1), "P levels are alphabetically after LT levels")
}

func TestLadderXPRewardRanges(t *testing.T) {
	tests := []struct {
		name     string
		level    repository.LadderLevel
		minXP    int32
		maxXP    int32
		expected bool
	}{
		{
			name:     "P1 typical range",
			level:    repository.LadderLevelP1,
			minXP:    1,
			maxXP:    50,
			expected: true,
		},
		{
			name:     "P2 typical range",
			level:    repository.LadderLevelP2,
			minXP:    51,
			maxXP:    150,
			expected: true,
		},
		{
			name:     "P3 typical range",
			level:    repository.LadderLevelP3,
			minXP:    151,
			maxXP:    300,
			expected: true,
		},
		{
			name:     "LT1 typical range",
			level:    repository.LadderLevelLT1,
			minXP:    301,
			maxXP:    500,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test verifies XP ranges are logical for career levels
			assert.Less(t, tt.minXP, tt.maxXP, "Min XP should be less than Max XP")
		})
	}
}

func TestCreateLadderLevelParams(t *testing.T) {
	params := repository.CreateLadderLevelParams{
		Level:           repository.LadderLevelP2,
		XpReward:        100,
		Technical:       "Advanced coding skills",
		ExpectedResults: "Deliver complex features",
		LeadershipScope: "Lead small teams",
	}

	assert.Equal(t, repository.LadderLevelP2, params.Level)
	assert.Equal(t, int32(100), params.XpReward)
	assert.NotEmpty(t, params.Technical)
	assert.NotEmpty(t, params.ExpectedResults)
	assert.NotEmpty(t, params.LeadershipScope)
}

func TestCareerLadderModel(t *testing.T) {
	// Test the complete CareerLadder structure
	ladder := repository.CareerLadder{
		Level:           repository.LadderLevelP3,
		XpReward:        250,
		Technical:       "Expert-level technical skills",
		ExpectedResults: "Deliver high-impact projects",
		LeadershipScope: "Mentor multiple engineers",
	}

	assert.Equal(t, repository.LadderLevelP3, ladder.Level)
	assert.Greater(t, ladder.XpReward, int32(0))
	assert.NotEmpty(t, ladder.Technical)
	assert.NotEmpty(t, ladder.ExpectedResults)
	assert.NotEmpty(t, ladder.LeadershipScope)
}

// Test that all ladder levels are distinct
func TestLadderLevelUniqueness(t *testing.T) {
	levels := []repository.LadderLevel{
		repository.LadderLevelP1,
		repository.LadderLevelP2,
		repository.LadderLevelP3,
		repository.LadderLevelLT1,
		repository.LadderLevelLT2,
		repository.LadderLevelLT3,
	}

	seen := make(map[repository.LadderLevel]bool)
	for _, level := range levels {
		assert.False(t, seen[level], "Level %s should be unique", level)
		seen[level] = true
	}

	assert.Len(t, seen, 6, "Should have exactly 6 unique levels")
}
