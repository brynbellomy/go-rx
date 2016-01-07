package rx

import "github.com/brynbellomy/go-result"

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
			out, cancel := NewSubject(), NewCancelable()

			go func() {
				defer cancelIn.Cancel()
				defer out.Complete()

				for {
					select {
					case <-cancel.OnCancel():
						return

					case x, open := <-chIn.Out():
						if !open {
							return
						}
						if x.IsError() {
							out.Send(x)
						} else {
							newVal, err := tfm(x.Value())
							if err != nil {
								out.Send(result.Failure(err))
							} else {
								out.Send(result.Success(newVal))
							}
						}
					}
				}
			}()

			return out, cancel
		},
	}
}

func (o *Observable) Subscribe() (IObservable, IDisposable) {
	return o.OnSubscribe()
}
