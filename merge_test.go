package rx_test

import (
	"sort"

	"github.com/brynbellomy/go-result"
	"github.com/brynbellomy/go-rx"

	τ "gopkg.in/check.v1"
)

type mergeSuite struct{}

var _ = τ.Suite(&mergeSuite{})

func (s *mergeSuite) TestMerge(c *τ.C) {
	ch1, ch2 := rx.NewSubject(), rx.NewSubject()

	m := rx.NewMerge(ch1, ch2)

	go func() {
		ch1.Send(result.Success("xyzzy"))
		ch1.Send(result.Success(1337))
		ch1.Send(result.Success("yeah"))
		ch1.Complete()
	}()

	go func() {
		ch2.Send(result.Success("foobar"))
		ch2.Send(result.Success(123))
		ch2.Send(result.Success(456))
		ch2.Complete()
	}()

	var (
		strings         []string
		ints            []int
		expectedStrings = []string{"xyzzy", "yeah", "foobar"}
		expectedInts    = []int{1337, 123, 456}
	)

	out, _ := m.Subscribe()

	for res := range out.Out() {
		c.Assert(res.IsError(), τ.Equals, false)

		x := res.Value()
		switch x := x.(type) {
		case string:
			strings = append(strings, x)
		case int:
			ints = append(ints, x)
		}
	}

	sort.Strings(strings)
	sort.Strings(expectedStrings)
	sort.Ints(ints)
	sort.Ints(expectedInts)

	c.Assert(strings, τ.HasLen, 3)
	c.Assert(ints, τ.HasLen, 3)

	for i := 0; i < 3; i++ {
		c.Assert(strings[i], τ.Equals, expectedStrings[i])
		c.Assert(ints[i], τ.Equals, expectedInts[i])
	}
}
