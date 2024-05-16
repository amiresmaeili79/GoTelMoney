package models

import "github.com/amir79esmaeili/go-tel-money/internal/messages"

type Conversation struct {
	Type  messages.ConvType
	State int
	Data  interface{}
}
