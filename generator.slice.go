package rx

import "github.com/brynbellomy/go-result"

type SliceGenerator struct {
	Slice []interface{}
}

func (g *SliceGenerator) Subscribe() (IObservable, IDisposable) {
	out, cancel := NewSubject(), NewCancelable()

	go func() {
		defer out.Complete()

		for _, x := range g.Slice {
			select {
			case <-cancel.OnCancel():
				return

			default:
				out.Send(result.Success(x))
			}
		}
	}()

	return out, cancel
}

func (g *SliceGenerator) AsObservable() *Observable {
	return &Observable{func() (IObservable, IDisposable) {
		return g.Subscribe()
	}}
}
