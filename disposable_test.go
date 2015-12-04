package rx_test

import (
	"github.com/brynbellomy/go-rx"
	τ "gopkg.in/check.v1"
)

type disposableSuite struct{}

var _ = τ.Suite(&disposableSuite{})

func (s *disposableSuite) TestFuncDisposable(c *τ.C) {
	wasCalled := false
	f := rx.NewFuncDisposable(func() { wasCalled = true })

	f.Cancel()

	c.Assert(wasCalled, τ.Equals, true)

	c.Assert(func() { f.Cancel() }, τ.Panics, "FuncDisposable: already disposed")
}

func (s *disposableSuite) TestCompositeDisposable(c *τ.C) {
	d := rx.NewCompositeDisposable()

	flag1, flag2, flag3 := false, false, false

	d.Add(rx.NewFuncDisposable(func() { flag1 = true }))
	d.Add(rx.NewFuncDisposable(func() { flag2 = true }))
	d.Add(rx.NewFuncDisposable(func() { flag3 = true }))

	d.Cancel()

	for _, flag := range []bool{flag1, flag2, flag3} {
		c.Assert(flag, τ.Equals, true)
	}

	c.Assert(func() { d.Cancel() }, τ.Panics, "CompositeDisposable: already disposed")
}
