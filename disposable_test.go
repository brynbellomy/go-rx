package async

import τ "gopkg.in/check.v1"

type disposableSuite struct{}

var _ = τ.Suite(&disposableSuite{})

func (s *disposableSuite) TestFuncDisposable(c *τ.C) {
	wasCalled := false
	f := NewFuncDisposable(func() { wasCalled = true })

	f.Cancel()

	c.Assert(wasCalled, τ.Equals, true)

	c.Assert(func() { f.Cancel() }, τ.Panics, "FuncDisposable: already disposed")
}

func (s *disposableSuite) TestCompositeDisposable(c *τ.C) {
	d := NewCompositeDisposable()

	flag1, flag2, flag3 := false, false, false

	d.Add(NewFuncDisposable(func() { flag1 = true }))
	d.Add(NewFuncDisposable(func() { flag2 = true }))
	d.Add(NewFuncDisposable(func() { flag3 = true }))

	d.Cancel()

	for _, flag := range []bool{flag1, flag2, flag3} {
		c.Assert(flag, τ.Equals, true)
	}

	c.Assert(func() { d.Cancel() }, τ.Panics, "CompositeDisposable: already disposed")
}
