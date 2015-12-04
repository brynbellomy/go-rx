package rx_test

import (
	"github.com/brynbellomy/go-rx"
	"github.com/listenonrepeat/listenonrepeat/backend/common/result"

	τ "gopkg.in/check.v1"
)

type concatSuite struct{}

var _ = τ.Suite(&concatSuite{})

func (s *concatSuite) TestConcat(c *τ.C) {
	ch1, ch2, ch3 := rx.NewSubject(), rx.NewSubject(), rx.NewSubject()
	cat := rx.NewConcat(ch1, ch2, ch3)

	go func() {
		ch1.Send(result.Success("xyzzy"))
		ch1.Send(result.Success("l33t"))
		ch1.Complete()
	}()

	go func() {
		ch2.Send(result.Success("yeah"))
		ch2.Complete()
	}()

	go func() {
		ch3.Send(result.Success("foobar"))
		ch3.Send(result.Success("zork"))
		ch3.Complete()
	}()

	recvd := make([]string, 0)
	expected := []string{"xyzzy", "l33t", "yeah", "foobar", "zork"}

	out, _ := cat.Subscribe()
	for res := range out.Out() {
		c.Assert(res.IsError(), τ.Equals, false)

		x := res.Value().(string)
		recvd = append(recvd, x)
	}

	c.Assert(recvd, τ.HasLen, len(expected))
	for i := 0; i < len(recvd); i++ {
		c.Assert(recvd[i], τ.Equals, expected[i])
	}
}
