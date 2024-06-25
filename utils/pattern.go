package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type Pattern struct {
	Name     string
	Quantity int
}

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

func GetSelectedPatternItemText(selectedIndex int, patterns []Pattern) string {
	selectedIndex = min(selectedIndex, len(patterns)-1)
	text := ""
	for i, pattern := range patterns {
		if i == selectedIndex {
			text += fmt.Sprintf(" > %s:%d", pattern.Name, pattern.Quantity)
		} else {
			text += fmt.Sprintf("   %s:%d", pattern.Name, pattern.Quantity)
		}
	}

	return text
}

func GetPatternItemText(patterns []Pattern) string {
	text := ""
	for _, pattern := range patterns {
		text += fmt.Sprintf("   %s:%d", pattern.Name, pattern.Quantity)
	}

	return text
}
