package rx

type Take struct {
	N int

	IObserver
	in *Subject
}

func NewTake(n int) *Take {
	in := NewSubject()

	return &Take{
		N:         n,
		in:        in,
		IObserver: in,
	}
}

func (t *Take) Subscribe() (IObservable, IDisposable) {
	out, cancel := NewSubject(), NewCancelable()

	go func() {
		defer out.Complete()

		i := 0
		for {
			select {
			case <-cancel.OnCancel():
				return

			case x, open := <-t.in.Out():
				if !open {
					return
				}

				if i >= t.N {
					return
				} else {
					i++
					out.Send(x)
				}
			}
		}
	}()

	return out, cancel
}
