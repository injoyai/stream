package pkg

import (
	"io"
	"sync"
)

type Input interface {
	io.Reader
}

type Output interface {
	io.Writer
}

func Copy(i Input, o Output) error {
	_, err := io.Copy(o, i)
	return err
}

type Inputer struct {
	Input
}

func (this *Inputer) WriteTo(w io.Writer) (int64, error) {
	return io.Copy(w, this.Input)
}

type Inputs struct {
	m  map[string]Input
	mu sync.RWMutex
}

func (this *Inputs) Get(key string) (Input, bool) {
	this.mu.RLock()
	v, ok := this.m[key]
	this.mu.RUnlock()
	return v, ok
}
