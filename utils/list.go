package utils

import (
	"github.com/charmbracelet/bubbles/list"
)

const (
	listHeightMargin = 5
	listWidthMargin  = 2
)

func NewList(items []list.Item, delegate list.ItemDelegate, name string) list.Model {
	width, height := GetWindowSize()
	width, height = getListSize(len(items), width, height)

	LogInfo("initializing %s list with size %d x %d", name, width, height)
	list := list.New(items, delegate, width, height)
	list.SetShowHelp(false)
	list.SetShowTitle(false)
	list.FilterInput.Cursor.Style = CursorStyle
	list.FilterInput.PromptStyle = ActionStyle
	list.DisableQuitKeybindings()
	list.ResetFilter()
	if len(items) == 0 {
		list.SetShowStatusBar(false)
	}

	return list
}

func SetListSize(list *list.Model, windowWidth, windowHeight int) {
	width, height := getListSize(len(list.Items()), windowWidth, windowHeight)
	list.SetSize(width, height)
	list.FilterInput.Width = 0
}

func getListSize(itemCount, windowWidth, windowHeight int) (int, int) {
	return windowWidth - listWidthMargin*2, getListHeight(itemCount, windowHeight)
}

func getListHeight(itemCount, windowHeight int) int {
	if itemCount == 0 {
		return 1
	}
	listMaxSize := itemCount*ListItemHeight + listHeightMargin
	return min(listMaxSize, windowHeight-listHeightMargin*2-TotalHelpWidth)
}
