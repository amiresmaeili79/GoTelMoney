package models

import (
	"github.com/amir79esmaeili/go-tel-money/internal/commands"
)

type Conversation struct {
	Type  commands.ConvType
	State int
	Data  interface{}
}

func NewConversation(t commands.ConvType, data interface{}) *Conversation {
	return &Conversation{
		Type:  t,
		State: 0,
		Data:  data,
	}
}
