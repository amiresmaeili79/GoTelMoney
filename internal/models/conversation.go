package models

import (
	"github.com/amir79esmaeili/go-tel-money/internal/commands"
)

type Conversation struct {
	Type  commands.ConvType
	State int
	Data  interface{}
}
