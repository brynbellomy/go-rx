package rx_test

import (
	"sync"

	"github.com/brynbellomy/go-rx"
	"github.com/listenonrepeat/listenonrepeat/backend/common/result"

	τ "gopkg.in/check.v1"
)

type batcherSuite struct{}

var _ = τ.Suite(&batcherSuite{})

func (s *batcherSuite) TestBatcher(c *τ.C) {
	b := rx.NewBatcher(5)

	wait := &sync.WaitGroup{}
	wait.Add(1)

	go func() {
		defer wait.Done()
		defer b.Complete()

		for i := 0; i < 12; i++ {
			b.Send(result.Success(i))
		}
	}()

	out, _ := b.Subscribe()

	x := <-out.Out()
	results, ok := x.Value().([]result.Result)
	c.Assert(ok, τ.Equals, true)
	c.Assert(results, τ.HasLen, 5)

	x = <-out.Out()
	results, ok = x.Value().([]result.Result)
	c.Assert(ok, τ.Equals, true)
	c.Assert(results, τ.HasLen, 5)

	x = <-out.Out()
	results, ok = x.Value().([]result.Result)
	c.Assert(ok, τ.Equals, true)
	c.Assert(results, τ.HasLen, 2)

	wait.Wait()
}
