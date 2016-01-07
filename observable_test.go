package rx_test

import (
	"sync"

	"github.com/brynbellomy/go-result"
	"github.com/brynbellomy/go-rx"
	τ "gopkg.in/check.v1"
)

type observableSuite struct{}

var _ = τ.Suite(&observableSuite{})

func (s *observableSuite) TestBroadcastObservableSliceGenerator(c *τ.C) {
	const BATCH_SIZE = 3

	sl := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9}
	gen := &rx.SliceGenerator{sl}

	b := gen.AsObservable().Broadcast()

	out1, _ := b.Subscribe()
	out2, _ := b.Subscribe()

	subscribers := []rx.IObservable{out1, out2}

	wait := sync.WaitGroup{}
	wait.Add(len(subscribers))

	for _, out := range subscribers {
		go func(out rx.IObservable) {
			i := 0
			for x := range out.Out() {
				c.Assert(x.Value(), τ.Equals, sl[i])
				i++
			}
			c.Assert(i, τ.Equals, len(sl))
			wait.Done()
		}(out)
	}
	wait.Wait()
}

func (s *observableSuite) TestBatchObservableSliceGenerator(c *τ.C) {
	const BATCH_SIZE = 3

	sl := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9}
	gen := &rx.SliceGenerator{sl}

	out, _ := gen.AsObservable().Batch(BATCH_SIZE).Subscribe()

	i := 0
	for batch := range out.Out() {
		xs := batch.Value().([]result.Result)
		c.Assert(xs, τ.HasLen, BATCH_SIZE)

		for _, x := range xs {
			c.Assert(x.Value(), τ.DeepEquals, sl[i])
			i++
		}
	}
}

func (s *observableSuite) TestMapObservableSliceGenerator(c *τ.C) {
	const BATCH_SIZE = 3

	sl := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9}
	gen := &rx.SliceGenerator{sl}

	out, _ := gen.AsObservable().Map(func(x interface{}) (interface{}, error) {
		return x.(int) * 2, nil
	}).Subscribe()

	i := 0
	for x := range out.Out() {
		c.Assert(x.Value(), τ.DeepEquals, sl[i].(int)*2)
		i++
	}

}
