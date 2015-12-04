package rx_test

import (
	"sort"

	"github.com/brynbellomy/go-rx"
	"github.com/listenonrepeat/listenonrepeat/backend/common/result"

	τ "gopkg.in/check.v1"
)

type mergeSuite struct{}

var _ = τ.Suite(&mergeSuite{})

func (s *mergeSuite) TestMerge(c *τ.C) {
	ch1, ch2 := make(chan result.Result), make(chan result.Result)

	m := rx.NewMerge(ch1, ch2)

	go func() {
		ch1 <- result.Success("xyzzy")
		ch1 <- result.Success(1337)
		ch1 <- result.Success("yeah")
		close(ch1)
	}()

	go func() {
		ch2 <- result.Success("foobar")
		ch2 <- result.Success(123)
		ch2 <- result.Success(456)
		close(ch2)
	}()

	var (
		strings         []string
		ints            []int
		expectedStrings = []string{"xyzzy", "yeah", "foobar"}
		expectedInts    = []int{1337, 123, 456}
	)

	m.Start()

	for res := range m.Out() {
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
