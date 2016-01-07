package rx_test

import (
	"github.com/brynbellomy/go-rx"
	τ "gopkg.in/check.v1"
)

type sliceGeneratorSuite struct{}

var _ = τ.Suite(&sliceGeneratorSuite{})

func (s *sliceGeneratorSuite) TestSliceGenerator(c *τ.C) {
	sl := []interface{}{1, 3, 5, 7, 9}

	g := &rx.SliceGenerator{sl}
	out, _ := g.Subscribe()
	i := 0
	for x := range out.Out() {
		c.Assert(x.Value(), τ.DeepEquals, sl[i])
		i++
	}
}
