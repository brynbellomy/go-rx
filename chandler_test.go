package async

import (
	"github.com/listenonrepeat/listenonrepeat/backend/common/result"
	τ "gopkg.in/check.v1"
)

type chanObserverSuite struct{}

var _ = τ.Suite(&chanObserverSuite{})

func (s *chanObserverSuite) TestChandler(c *τ.C) {
	ch := make(chan result.Result)
	chClosed := make(chan struct{})
	recvd := make([]result.Result, 0)
	length := 0

	b := NewChandler(ch, Handlers{
		OnNext: func(r result.Result) {
			recvd = append(recvd, r)
		},
		OnComplete: func() {
			length = len(recvd)
			close(chClosed)
		},
	})
	b.Start()

	go func() {
		ch <- result.Success(10)
		ch <- result.Success(20)
		ch <- result.Success(30)
		close(ch)
	}()

	<-chClosed

	c.Assert(length, τ.Equals, 3)
	c.Assert(recvd, τ.HasLen, 3)

	c.Assert(recvd[0].Value(), τ.Equals, 10)
	c.Assert(recvd[1].Value(), τ.Equals, 20)
	c.Assert(recvd[2].Value(), τ.Equals, 30)

	c.Assert(recvd[0].IsError(), τ.Equals, false)
	c.Assert(recvd[1].IsError(), τ.Equals, false)
	c.Assert(recvd[2].IsError(), τ.Equals, false)
}
