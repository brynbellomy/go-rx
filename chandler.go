package async

import "github.com/listenonrepeat/listenonrepeat/backend/common/result"

type (
	Chandler struct {
		Chan <-chan result.Result
		Handlers

		*Cancelable
	}

	Handlers struct {
		OnNext     func(x result.Result)
		OnComplete func()
	}
)

func NewChandler(ch <-chan result.Result, handlers Handlers) *Chandler {
	if ch == nil {
		panic("Chandler: chan is nil")
	} else if handlers.OnNext == nil {
		panic("Chandler: Handlers.OnNext is nil")
	}

	return &Chandler{
		Chan:       ch,
		Handlers:   handlers,
		Cancelable: NewCancelable(),
	}
}

func (c *Chandler) Start() {
	go func() {
		defer c.dispose()

	Outer:
		for {
			select {
			case <-c.OnCancel():
				break Outer

			case r, open := <-c.Chan:
				if open {
					c.OnNext(r)
				} else {
					// channel was closed, so we shut down the subscriber
					break Outer
				}
			}
		}
	}()
}

func (c *Chandler) dispose() {
	if c.OnComplete != nil {
		c.OnComplete()
	}
}
