package main

import (
	"bytes"
	"testing"
)

func stubCommand(x int) RawCommand{
	return []byte{byte(x)}
}

const N = 20

func TestHistoryCreateChan_IterateAllCommand(t *testing.T) {
	h := CreateHistory()
	ch := h.CreateChan(0)
	for i := 0; i < N/2; i++ {
		h.AppendCommand(stubCommand(i))
	}
	go func() {
		for i := N/2; i < N; i++ {
			h.AppendCommand(stubCommand(i))
		}
	}()

	for i := 0; i < N; i++ {
		if(!bytes.Equal(<-ch, stubCommand(i))) {
			t.Fatalf("%d'th Command is not expected: %s", i, <-ch)
		}
	}

}

func TestHistory_CanCopy(t *testing.T) {
	h1 := CreateHistory()
	ch := h1.CreateChan(0)
	h2 := h1
	for i := 0; i < N; i++ {
		if i|1 == 0 {
			h1.AppendCommand(stubCommand(i))
		} else {
			h2.AppendCommand(stubCommand(i))
		}
	}
	
	for i := 0; i < N; i++ {
		if(!bytes.Equal(<-ch, stubCommand(i))) {
			t.Fatalf("%d'th Command is not expected: %s", i, <-ch)
		}
	}
	

}
