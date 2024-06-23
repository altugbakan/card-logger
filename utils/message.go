package utils

import "fmt"

type Message interface {
	Render() string
}

type InfoMessage struct {
	Message
	Text string
}

func (m InfoMessage) Render() string {
	return DimTextStyle.Render(m.Text)
}

func NewInfoMessage(format string, a ...any) InfoMessage {
	return InfoMessage{Text: fmt.Sprintf(format, a...)}
}

type ErrorMessage struct {
	Message
	Text string
}

func (m ErrorMessage) Render() string {
	return ErrorStyle.Render(m.Text)
}

func NewErrorMessage(format string, a ...any) ErrorMessage {
	return ErrorMessage{Text: fmt.Sprintf(format, a...)}
}
