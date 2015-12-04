package async

import (
	"sync"

	τ "gopkg.in/check.v1"
)

type cancelableSuite struct{}

var _ = τ.Suite(&cancelableSuite{})

func (s *cancelableSuite) TestCancelable(c *τ.C) {
	b := NewCancelable()

	canceled := []bool{false, false, false}

	wait := sync.WaitGroup{}
	wait.Add(3)
	for i := 0; i < 3; i++ {
		go func(i int) {
			<-b.OnCancel()
			canceled[i] = true
			wait.Done()
		}(i)
	}

	b.Cancel()

	wait.Wait()
	c.Assert(canceled[0], τ.Equals, true)
	c.Assert(canceled[1], τ.Equals, true)
	c.Assert(canceled[2], τ.Equals, true)
}
