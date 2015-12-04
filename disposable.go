package rx

type (
	FuncDisposable struct {
		onCancel func()
		disposed bool
	}

	CompositeDisposable struct {
		disposables []IDisposable
		disposed    bool
	}
)

func NewFuncDisposable(onCancel func()) *FuncDisposable {
	return &FuncDisposable{onCancel: onCancel}
}

func (f *FuncDisposable) Cancel() {
	if f.disposed {
		panic("FuncDisposable: already disposed")
	}
	f.disposed = true

	f.onCancel()
}

func NewCompositeDisposable(disposables ...IDisposable) *CompositeDisposable {
	return &CompositeDisposable{disposables: disposables}
}

func (c *CompositeDisposable) Add(others ...IDisposable) {
	for _, d := range others {
		c.disposables = append(c.disposables, d)
	}
}

func (c *CompositeDisposable) Cancel() {
	if c.disposed {
		panic("CompositeDisposable: already disposed")
	}
	c.disposed = true

	for _, d := range c.disposables {
		d.Cancel()
	}

	c.disposables = nil
}
