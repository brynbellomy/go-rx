package rx

import "sync"

type Merge struct {
	Inputs []ISubscribable
}

func NewMerge(inputs ...ISubscribable) *Merge {
	if inputs == nil || len(inputs) == 0 {
		panic("Merge requires at least one input")
	}

	return &Merge{Inputs: inputs}
}

func (m *Merge) Subscribe() (IObservable, IDisposable) {
	dispose := NewCompositeDisposable()

	out, cancel := NewSubject(), NewCancelable()
	dispose.Add(cancel)

	wait := sync.WaitGroup{}
	wait.Add(len(m.Inputs))

	for _, in := range m.Inputs {
		chIn, cancel := in.Subscribe()
		dispose.Add(cancel)

		go func() {
			defer wait.Done()

			for x := range chIn.Out() {
				out.Send(x)
			}
		}()
	}

	// Start a goroutine to close `out` once all the output goroutines are
	// done.  This must start after the wait.Add call.
	go func() {
		defer out.Complete()
		wait.Wait()
	}()

	return out, cancel
}
