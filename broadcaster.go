package async

import "github.com/listenonrepeat/listenonrepeat/backend/common/result"

type Broadcaster struct {
	channels  map[uint64]chan result.Result
	completed bool
	counter   uint64
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		channels: make(map[uint64]chan result.Result),
	}
}

func (b *Broadcaster) Subscribe() (IObservable, IDisposable) {
	ch := make(chan result.Result)

	counter := b.counter
	b.counter++

	b.channels[counter] = ch

	return ObservableChan(ch), NewFuncDisposable(func() {
		close(ch)
		delete(b.channels, counter)
	})
}

func (b *Broadcaster) Send(val result.Result) {
	for _, ch := range b.channels {
		ch <- val
	}
}

func (b *Broadcaster) Complete() {
	if b.completed {
		panic("Broadcaster: cannot call .Complete() more than once")
	}

	for _, ch := range b.channels {
		close(ch)
	}
	b.channels = nil
	b.completed = true
}
