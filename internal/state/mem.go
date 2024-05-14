package state

import (
	"fmt"
)

type InMemoryState struct {
	State map[int64]interface{}
}

func NewInMemoryState() *InMemoryState {
	return &InMemoryState{
		State: make(map[int64]interface{}),
	}
}

func (mem *InMemoryState) CreateSession(chatId int64, data interface{}) {
	mem.State[chatId] = data
}

func (mem *InMemoryState) GetSession(chatId int64) (interface{}, error) {
	conv, exists := mem.State[chatId]
	if !exists {
		return nil, fmt.Errorf("session %d does not exist", chatId)
	}
	return conv, nil
}

func (mem *InMemoryState) ChangeState(chatId int64, state string) error {
	_, exists := mem.State[chatId]
	if !exists {
		return fmt.Errorf("session %d does not exist", chatId)
	}
	return nil
}

func (mem *InMemoryState) DeleteState(chatId int64) error {
	_, exists := mem.State[chatId]
	if !exists {
		return fmt.Errorf("session %d does not exist", chatId)
	}
	delete(mem.State, chatId)
	return nil
}
