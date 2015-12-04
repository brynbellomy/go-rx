package rx

import "github.com/listenonrepeat/listenonrepeat/backend/common/result"

type Batcher struct {
	size int

	IObserver
	IObservable

	in  *Subject // buffer size of `in` is determined by the value of `.size`
	out *Subject
}

func NewBatcher(size int) *Batcher {
	if size <= 0 {
		panic("NewBatcher: size must be a positive, non-zero integer")
	}

	in := NewBufferedSubject(size)
	out := NewSubject()

	return &Batcher{
		size:        size,
		IObserver:   in,
		IObservable: out,
		in:          in,
		out:         out,
	}
}

func (b *Batcher) Subscribe() (IObservable, IDisposable) {
	cancel := NewCancelable()

	go func() {
		defer b.dispose()

		i := 0
		current := make([]result.Result, 0, b.size)

	OuterLoop:
		for {
			select {
			case <-cancel.OnCancel():
				return

			case x, open := <-b.in.Out():
				if !open {
					break OuterLoop
				}

				current = append(current, x)
				i++
				if i == b.size {
					b.sendBatch(current)
					i = 0
					current = make([]result.Result, 0, b.size)
				}
			}
		}

		if len(current) > 0 {
			b.sendBatch(current)
		}
	}()

	return b.out, cancel
}

func (b *Batcher) dispose() {
	b.out.Complete()
}

func (b *Batcher) sendBatch(batch []result.Result) {
	outBatch := make([]result.Result, len(batch))
	copy(outBatch, batch)
	b.out.Send(result.Success(outBatch))
}
