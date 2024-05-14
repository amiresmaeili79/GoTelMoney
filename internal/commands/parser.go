package commands

import (
	"fmt"
	"strings"

	"github.com/amir79esmaeili/go-tel-money/internal/conversations"
)

func WhatIsCommand(msg string) (conversations.ConvType, error) {
	for key, value := range conversations.Commands {
		if strings.Contains(msg, key) {
			return value, nil
		}
	}
	return -1, fmt.Errorf("unknown command %s", msg)
}
