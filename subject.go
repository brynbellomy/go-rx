package rx

import "github.com/brynbellomy/go-result"

type Subject struct {
	ch        chan result.Result
	completed bool
	DebugName string
}

func NewSubject() *Subject {
	return &Subject{ch: make(chan result.Result)}
}

func NewBufferedSubject(size int) *Subject {
	return &Subject{ch: make(chan result.Result, size)}
}

func (s *Subject) Send(r result.Result) {
	if s.completed {
		panic("Subject: cannot .Send() after completing")
	}
	s.ch <- r
}

func (s *Subject) Complete() {
	if s.completed {
		panic("Subject: cannot .Complete() more than once")
	}
	close(s.ch)
}

func (s *Subject) Out() ObservableChan {
	return s.ch
}

func (s *Subject) Subscribe() (IObservable, IDisposable) {
	return s.Out(), NewFuncDisposable(s.Complete)
}
