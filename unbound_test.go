package main

import (
	"testing"
	"sync"
	"time"
)

func TestMakeInfiniteNoPause(t *testing.T) {
	in, out := Makelnfinite2()
	lastVal := -1
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for v := range out {
			vi := v.(int)
			if lastVal+1 != vi {
				t.Errorf("Unexpected value; expected %d, got %d", lastVal+1, vi)
			}
			lastVal = vi
		}
		wg.Done()
	}()
	for i:=0;i<100;i++ {
		time.Sleep(1*time.Millisecond)
		in<-i
	}
	close(in)
	wg.Wait()
	if lastVal != 99 {
		t.Errorf("Didn't get all values, last one received was %d", lastVal)
	}
}

func Makelnfinite2() (chan<- interface{}, <-chan interface{}) {
	in := make(chan interface{})
	out := make(chan interface{})
	im := &sync.RWMutex{}
	im.Lock()
	var inQueue []interface{}
	index := 0

	go func() {
		for in != nil {
			v, ok := <- in
			if(!ok) {
				in = nil
				im.Unlock()
			} else {
				inQueue = append(inQueue, v)
				im.Unlock()
				im = &sync.RWMutex{}
				im.Lock()
			}
		}
	}()

	go func() {
		for {
			if index < len(inQueue) {
				out <- inQueue[index]
				index++
			} else {
				im.RLock()
				if index < len(inQueue) {
					out <- inQueue[index]
					index++
				} else {
					close(out)
					break;
				}
			}
		}
	}()
	return in, out
}

func Makelnfinite() (chan<- interface{}, <-chan interface{}) {
	in := make(chan interface{})
	out := make(chan interface{})
	go func() {
		var inQueue []interface{}
		outCh := func() chan interface{} {
			if len(inQueue) == 0 {
				return nil
			}
			return out
		}
		curVal := func() interface{} {
			if len(inQueue) == 0 {
				return nil
			}
			return inQueue[0]
		}
		for len(inQueue) > 0 || in != nil {
			select {
			case v, ok := <-in:
				if !ok {
					in = nil
				} else {
					inQueue = append(inQueue, v)
				}
			case outCh() <- curVal():
				inQueue = inQueue[1:]
			}
		}
		close(out)
	}()
	return in, out
}

