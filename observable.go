package rx

type Observable struct {
	OnSubscribe func() (IObservable, IDisposable)
}

func (o *Observable) Batch(size int) *Observable {
	return &Observable{func() (IObservable, IDisposable) {
		chIn, cancelIn := o.Subscribe()
		b := NewBatcher(size)
		go func() {
			defer cancelIn.Cancel()
			for x := range chIn.Out() {
				b.Send(x)
			}
			b.Complete()
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
