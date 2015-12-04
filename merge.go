package rx

import (
	"sync"

	"github.com/listenonrepeat/listenonrepeat/backend/common/result"
)

type Merge struct {
	Inputs []<-chan result.Result

	*Cancelable

	out  *Subject
	wait sync.WaitGroup

	IObservable
}

func NewMerge(inputs ...<-chan result.Result) *Merge {
	if inputs == nil || len(inputs) == 0 {
		panic("Merge requires at least one input")
	}

	out := NewSubject()

	return &Merge{
		Inputs:      inputs,
		Cancelable:  NewCancelable(),
		out:         out,
		IObservable: out,
	}
}

func (m *Merge) Start() {
	// Start an output goroutine for each input channel in m.Inputs.  output
	// copies values from c to out until c is closed, then calls wait.Done.
	m.wait.Add(len(m.Inputs))
	for _, c := range m.Inputs {
		go m.pipe(c)
	}

	// Start a goroutine to close `chOut` once all the output goroutines are
	// done.  This must start after the wait.Add call.
	go func() {
		defer m.out.Complete()
		m.wait.Wait()
	}()
}

func (m *Merge) pipe(ch <-chan result.Result) {
	defer m.wait.Done()
	for res := range ch {
		m.out.Send(res)
	}
}
