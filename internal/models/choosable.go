package models

type InlineKeyboardItem interface {
	ToDisplay() string
	StringID() string
}
