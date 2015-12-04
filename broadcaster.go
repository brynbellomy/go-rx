package rx

import "github.com/listenonrepeat/listenonrepeat/backend/common/result"

type Broadcaster struct {
	Channels  map[uint64]chan result.Result
	completed bool
	counter   uint64
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		Channels: make(map[uint64]chan result.Result),
	}
}

func (b *Broadcaster) Subscribe() (IObservable, IDisposable) {
	ch := make(chan result.Result)

	counter := b.counter
	b.counter++

	b.Channels[counter] = ch

	return ObservableChan(ch), NewFuncDisposable(func() {
		close(ch)
		delete(b.Channels, counter)
	})
}

func (b *Broadcaster) Send(val result.Result) {
	for _, ch := range b.Channels {
		ch <- val
	}
}

func (b *Broadcaster) Complete() {
	if b.completed {
		panic("Broadcaster: cannot call .Complete() more than once")
	}

	for _, ch := range b.Channels {
		close(ch)
	}
	b.Channels = nil
	b.completed = true
}
