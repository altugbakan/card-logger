package utils

import (
	"strconv"
	"strings"
)

func CorrectPattern(pattern string) string {
	pattern = strings.ToUpper(pattern)
	switch pattern {
	case "HF":
		return "H"
	case "RHF", "R":
		return "RH"
	default:
		return pattern
	}
}

func GetPatternText(rarity string, possiblePatterns []string, patterns map[string]int) string {
	patternsText := make([]string, 0, len(possiblePatterns))
	for _, pattern := range possiblePatterns {
		quantity := patterns[pattern]

		patternsText = append(patternsText, pattern+": "+strconv.Itoa(quantity))
	}

	return strings.Join(patternsText, ", ")
}
