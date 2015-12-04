package rx

type (
	Concat struct {
		streams []ISubscribable
	}
)

func NewConcat(streams ...ISubscribable) *Concat {
	return &Concat{streams}
}

func (c *Concat) Subscribe() (IObservable, IDisposable) {
	out, cancel := NewSubject(), NewCancelable()

	go func() {
		defer out.Complete()

		for _, stream := range c.streams {
			s, _ := stream.Subscribe()

			for r := range s.Out() {
				out.Send(r)
			}
		}
	}()

	return out, cancel
}
