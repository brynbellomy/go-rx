package rx

type Observable struct {
	OnSubscribe func() (IObservable, IDisposable)
}

func ObservableFrom(subscribable ISubscribable) *Observable {
	return &Observable{subscribable.Subscribe}
}

func (o *Observable) Take(n int) *Observable {
	return &Observable{func() (IObservable, IDisposable) {
		chIn, cancelIn := o.Subscribe()
		t := NewTake(n)

		go func() {
			defer cancelIn.Cancel()
			defer t.Complete()

			for x := range chIn.Out() {
				t.Send(x)
			}
		}()
		return t.Subscribe()
	}}
}

func (o *Observable) Broadcast() *Broadcaster {
	chIn, cancelIn := o.Subscribe()
	b := NewBroadcaster()

	go func() {
		defer cancelIn.Cancel()
		defer b.Complete()

		for x := range chIn.Out() {
			b.Send(x)
		}
	}()

	return b
}

func (o *Observable) Batch(size int) *Observable {
	return &Observable{func() (IObservable, IDisposable) {
		chIn, cancelIn := o.Subscribe()
		b := NewBatcher(size)

		go func() {
			defer cancelIn.Cancel()
			defer b.Complete()

			for x := range chIn.Out() {
				b.Send(x)
			}
		}()
		return b.Subscribe()
	}}
}

func (o *Observable) Map(tfm func(x interface{}) (interface{}, error)) *Observable {
	return &Observable{
		func() (IObservable, IDisposable) {
			chIn, cancelIn := o.Subscribe()
			m := NewMap(tfm)

			go func() {
				defer cancelIn.Cancel()

				for x := range chIn.Out() {
					m.Send(x)
				}
				m.Complete()
			}()
			return m.Subscribe()
		},
	}
}

func (o *Observable) Subscribe() (IObservable, IDisposable) {
	return o.OnSubscribe()
}
