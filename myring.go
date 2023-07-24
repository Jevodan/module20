package main

import "sync"

type Ring struct {
	next, prev *Ring
	Value      int
	RW         sync.RWMutex
}

func (r *Ring) init() *Ring {
	r.next = r
	r.prev = r
	return r
}

func (r *Ring) Next() *Ring {
	if r.next == nil {
		return r.init()
	}
	return r.next
}

func NewRing(n int) *Ring {
	if n <= 0 {
		return nil
	}
	r := new(Ring)
	p := r
	for i := 1; i <= n; i++ {
		p.next = &Ring{prev: p}
		p = p.next
	}
	p.next = r
	r.prev = p
	return r
}

func (r *Ring) Len() int {
	n := 0
	if r != nil {
		for p := r.Next(); p != r; p = p.next {
			n++
		}
	}
	return n
}

func (r *Ring) Get() []int {
	data := make([]int, 0)
	if r != nil {
		for p := r.Next(); p != r; p = p.next {
			if p != nil {
				data = append(data, p.Value)
				p.Value = 0
			}
		}
	}
	return data
}
