package ring

import "sync"

type Ring struct {
	next, prev *Ring
	Value      int
	RW         sync.RWMutex
}

type RingInterface interface {
	Next() *Ring
	Len() int
	Get() []int
	init() *Ring
	SetValue()
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
		r.RW.RLock()
		for p := r.Next(); p != r; p = p.next {
			if p != nil {
				data = append(data, p.Value)
				p.Value = 0
			}
		}
		r.RW.RUnlock()
	}
	return data
}

func (r *Ring) SetValue(val int) {
	r.RW.Lock()
	r.Value = val
	r.RW.Unlock()
}
