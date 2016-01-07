package rx

import "github.com/brynbellomy/go-result"

type Broadcaster struct {
	Channels  map[uint64]IObserver
	completed bool
	counter   uint64
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		Channels: make(map[uint64]IObserver),
	}
}

func (b *Broadcaster) Subscribe() (IObservable, IDisposable) {
	out := NewSubject()

	counter := b.counter
	b.counter++

	b.Channels[counter] = out

	return out, NewFuncDisposable(func() {
		out.Complete()
		delete(b.Channels, counter)
	})
}

func (b *Broadcaster) Send(val result.Result) {
	for _, ch := range b.Channels {
		ch.Send(val)
	}
}

func (b *Broadcaster) Complete() {
	if b.completed {
		panic("Broadcaster: cannot call .Complete() more than once")
	}

	for _, ch := range b.Channels {
		ch.Complete()
	}
	b.Channels = nil
	b.completed = true
}
