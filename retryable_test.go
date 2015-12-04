package rx_test

import (
	"fmt"

	"github.com/brynbellomy/go-rx"
	τ "gopkg.in/check.v1"
)

type retryableSuite struct{}

var _ = τ.Suite(&retryableSuite{})

func (s *retryableSuite) TestRetryableSucceedImmediately(c *τ.C) {
	i := 0
	r := rx.Retryable{
		MaxAttempts: 5,
		Operation: rx.FuncOperation(func() (interface{}, error) {
			i++
			return 1337, nil
		}),
	}

	_, err := r.Execute()

	c.Assert(err, τ.IsNil)
	c.Assert(i, τ.Equals, 1)
}

func (s *retryableSuite) TestRetryableSucceedEventually(c *τ.C) {
	i := 0
	theErr := fmt.Errorf("the error")

	r := rx.Retryable{
		MaxAttempts: 5,
		Operation: rx.FuncOperation(func() (interface{}, error) {
			i++
			if i == 2 {
				return 1337, nil
			} else {
				return nil, theErr
			}
		}),
	}

	_, err := r.Execute()

	c.Assert(err, τ.IsNil)
	c.Assert(i, τ.Equals, 2)
}

func (s *retryableSuite) TestRetryableFailAlways(c *τ.C) {
	i := 0
	theErr := fmt.Errorf("the error")
	r := rx.Retryable{
		MaxAttempts: 5,
		Operation: rx.FuncOperation(func() (interface{}, error) {
			i++
			return nil, theErr
		}),
	}

	_, err := r.Execute()

	c.Assert(err, τ.Equals, theErr)
	c.Assert(i, τ.Equals, 4)
	c.Assert(r.Attempts, τ.Equals, 4)
}
