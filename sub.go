package main

import "sync"

func NewSubscribe() *Subscribe {
	return &Subscribe{}
}

type Chan struct {
	C         chan []byte
	once      sync.Once
	closeFunc func()
}

func (this *Chan) Close() {
	this.once.Do(this.closeFunc)
}

type Subscribe struct {
	ls []*Chan
}

func (this *Subscribe) Publish(data []byte) {
	for _, v := range this.ls {
		select {
		case v.C <- data:
		default:
		}
	}
}

func (this *Subscribe) Subscribe(size int) *Chan {
	c := make(chan []byte, size)
	ch := &Chan{
		C: c,
		closeFunc: func() {
			for i, v := range this.ls {
				if v.C == c {
					this.ls = append(this.ls[:i], this.ls[i+1:]...)
					return
				}
			}
		},
	}
	this.ls = append(this.ls, ch)
	return ch
}
