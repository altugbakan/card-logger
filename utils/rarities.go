package utils

import (
	"strings"
)

var rarityToPatterns = map[string][]string{
	"Common":                    {"N", "RH"},
	"Uncommon":                  {"N", "RH"},
	"Rare":                      {"H", "RH"},
	"Double Rare":               {"H"},
	"ACE SPEC Rare":             {"H"},
	"Illustration Rare":         {"H"},
	"Ultra Rare":                {"H"},
	"Special Illustration Rare": {"H"},
	"Hyper Rare":                {"H"},
}

func CorrectPattern(pattern string) string {
	switch pattern {
	case "HF":
		return "H"
	case "RHF", "R":
		return "RH"
	default:
		return pattern
	}
}

func IsPatternValidForRarity(pattern, rarity string) bool {
	patterns, ok := rarityToPatterns[rarity]
	if !ok {
		return false
	}

	pattern = strings.ToUpper(pattern)

	for _, p := range patterns {
		if p == pattern {
			return true
		}
	}

	return false
}
