package utils

import "github.com/charmbracelet/x/term"

func GetWindowSize() (int, int) {
	width, height, err := term.GetSize(0)
	if err != nil {
		LogError("could not get terminal size, using defaults: %v", err)
		width = DefaultWidth
		height = DefaultHeight
	}
	return width, height
}
