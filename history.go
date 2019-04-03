package main

import (
	"sync"
)

/*
RawCommand ...
*/
type RawCommand []byte

type history struct {
	commands []RawCommand
	mutex    *sync.RWMutex
}

/*
History ...
*/
type History interface {
	AppendCommand(command RawCommand)
	CreateChan(index int) <-chan RawCommand
}

/*
CreateHistory ...
*/
func CreateHistory() History {
	h := history{
		mutex: &sync.RWMutex{},
	}
	h.mutex.Lock()
	return &h
}

func (h *history) AppendCommand(command RawCommand) {
	h.commands = append(h.commands, command)
	h.mutex.Unlock()
	h.mutex = &sync.RWMutex{}
	h.mutex.Lock()
}

func (h *history) CreateChan(index int) <-chan RawCommand {
	out := make(chan RawCommand)
	go func() {
		for {
			if index < len(h.commands) {
				out <- h.commands[index]
				index++
			} else {
				h.mutex.RLock()
				if index < len(h.commands) {
					out <- h.commands[index]
					index++
				} else {
					close(out)
					break
				}
			}
		}
	}()
	return out
}
