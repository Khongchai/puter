package lib

import "sync"

// Promise-like stucture for carrying a value in each line evaluation.
//
// When a line is evaluated, the result might not yet be ready. For example
//
// x = 2 btc in usd
//
// requires fetching exchange rate from bitcoin to usd first.
type Promise[T any] struct {
	once  sync.Once
	ch    chan struct{}
	value T
}

func NewPromise[T any]() *Promise[T] {
	return &Promise[T]{}
}

func NewResolvedPromise[T any](v T) *Promise[T] {
	p := &Promise[T]{}
	return p.Resolve(v)
}

func (p *Promise[T]) Resolve(value T) *Promise[T] {
	p.once.Do(func() {
		p.value = value
		close(p.ch)
	})
	return p
}

func (p *Promise[T]) Await() T {
	<-p.ch
	return p.value
}
