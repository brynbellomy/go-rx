package rx

import "github.com/listenonrepeat/listenonrepeat/backend/common/result"

type (
	IGenerator interface {
		Next() (r result.Result, done bool)
	}

	ISubscribable interface {
		Subscribe() (out IObservable, cancel IDisposable)
	}

	IDisposable interface {
		Cancel()
	}

	IObservable interface {
		Out() ObservableChan
	}

	IObserver interface {
		Send(r result.Result)
		Complete()
	}

	ObservableChan <-chan result.Result
)

// IObservable interface conformance
func (c ObservableChan) Out() ObservableChan {
	return c
}

// IGenerator interface conformance
func (c ObservableChan) Next() (r result.Result, done bool) {
	r, done = <-c
	return
}
