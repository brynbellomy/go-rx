package rx_test

import (
	"github.com/brynbellomy/go-result"
	"github.com/brynbellomy/go-rx"
	τ "gopkg.in/check.v1"
)

type takeSuite struct{}

var _ = τ.Suite(&takeSuite{})

func (s *takeSuite) TestTake(c *τ.C) {
	sl := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	t := rx.NewTake(5)

	go func() {
		for _, x := range sl {
			t.Send(result.Success(x))
		}
	}()

	out, _ := t.Subscribe()
	i := 0
	for x := range out.Out() {
		c.Assert(x.Value().(int), τ.Equals, sl[i])
		i++
	}

	c.Assert(i, τ.Equals, 5)
}
