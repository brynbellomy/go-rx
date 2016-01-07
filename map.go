package rx

import "github.com/brynbellomy/go-result"

type Map struct {
	Transform MapFunc

	IObserver
	IObservable

	in  *Subject // buffer size of `in` is determined by the value of `.size`
	out *Subject
}

type MapFunc func(x interface{}) (interface{}, error)

func NewMap(tfm MapFunc) *Map {
	in, out := NewSubject(), NewSubject()

	return &Map{
		Transform:   tfm,
		IObserver:   in,
		IObservable: out,
		in:          in,
		out:         out,
	}
}

func (m *Map) Subscribe() (IObservable, IDisposable) {
	cancel := NewCancelable()

	go func() {
		defer m.out.Complete()

		for {
			select {
			case <-cancel.OnCancel():
				return

			case x, open := <-m.in.Out():
				if !open {
					return
				}

				if x.IsError() {
					m.out.Send(x)
				} else {
					newVal, err := m.Transform(x.Value())
					if err != nil {
						m.out.Send(result.Failure(err))
					} else {
						m.out.Send(result.Success(newVal))
					}
				}
			}
		}
	}()

	return m.out, cancel
}
