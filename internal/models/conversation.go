package models

import "github.com/amir79esmaeili/go-tel-money/internal/conversations"

type Conversation struct {
	Type  conversations.ConvType
	State int
	Data  interface{}
}
