package async

import (
	"fmt"
	τ "gopkg.in/check.v1"
)

type operationSuite struct{}

var _ = τ.Suite(&operationSuite{})

func (s *operationSuite) TestFuncOperationFailure(c *τ.C) {
	called := false
	theErr := fmt.Errorf("An error")

	op := NewFuncOperation(func() (interface{}, error) {
		called = true
		return nil, theErr
	})

	val, err := op.Execute()

	c.Assert(called, τ.Equals, true)
	c.Assert(val, τ.IsNil)
	c.Assert(err, τ.Equals, theErr)
}

func (s *operationSuite) TestFuncOperationSuccess(c *τ.C) {
	called := false

	op := NewFuncOperation(func() (interface{}, error) {
		called = true
		return 1337, nil
	})

	val, err := op.Execute()

	c.Assert(called, τ.Equals, true)
	c.Assert(err, τ.IsNil)
	c.Assert(val, τ.Equals, 1337)
}
