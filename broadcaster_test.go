package async

import (
	"sync"

	"github.com/listenonrepeat/listenonrepeat/backend/common/result"

	τ "gopkg.in/check.v1"
)

type broadcasterSuite struct{}

var _ = τ.Suite(&broadcasterSuite{})

func (s *broadcasterSuite) TestBroadcaster(c *τ.C) {
	b := NewBroadcaster()
	wait := &sync.WaitGroup{}

	chans := make([]IObservable, 3)
	chans[0], _ = b.Subscribe()
	chans[1], _ = b.Subscribe()
	chans[2], _ = b.Subscribe()

	c.Assert(b.channels, τ.HasLen, len(chans))

	wait.Add(1)
	go func() {
		defer wait.Done()
		b.Send(result.Success(123))
		b.Send(result.Success(456))
		b.Send(result.Success(789))
	}()

	for i := 0; i < 3; i++ {
		wait.Add(1)
		go func(i int) {
			defer wait.Done()
			x := <-chans[i].Out()
			c.Assert(x.Value(), τ.Equals, 123)
			x = <-chans[i].Out()
			c.Assert(x.Value(), τ.Equals, 456)
			x = <-chans[i].Out()
			c.Assert(x.Value(), τ.Equals, 789)
		}(i)
	}

	wait.Wait()
}
