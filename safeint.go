package main

import "sync"

type safeInt struct {
	val int
	mux sync.Mutex
}

func (i *safeInt) Set(val int) {
	i.mux.Lock()
	defer i.mux.Unlock()
	i.val = val
}

func (i *safeInt) Get() int {
	i.mux.Lock()
	defer i.mux.Unlock()
	return i.val
}
