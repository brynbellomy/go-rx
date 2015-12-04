package rx

type (
	Cancelable struct {
		ch       chan struct{}
		canceled bool
	}
)

func NewCancelable() *Cancelable {
	return &Cancelable{ch: make(chan struct{})}
}

func (c *Cancelable) Cancel() {
	if c.canceled {
		panic("async.Cancelable: Cannot call .Cancel() more than once")
	}
	c.canceled = true

	close(c.ch)
}

func (c *Cancelable) OnCancel() <-chan struct{} {
	return c.ch
}
