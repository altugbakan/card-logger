package utils

type Message interface {
	Render() string
}

type InfoMessage struct {
	Text string
}

func (m InfoMessage) Render() string {
	return DimTextStyle.Render(m.Text)
}

func NewInfoMessage(text string) InfoMessage {
	return InfoMessage{Text: text}
}

type ErrorMessage struct {
	Text string
}

func (m ErrorMessage) Render() string {
	return ErrorStyle.Render(m.Text)
}

func NewErrorMessage(text string) ErrorMessage {
	return ErrorMessage{Text: text}
}
