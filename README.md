
# go-rx

[![Build Status](https://drone.io/github.com/brynbellomy/go-rx/status.png)](https://drone.io/github.com/brynbellomy/go-rx/latest)

`go-rx` is an encapsulation of some basic, common patterns using Go's concurrency primitives.  It draws *very loosely* from the metaphor/architecture used by Microsoft's Reactive Extensions (Rx).

I've heard people say that Rx and Go shouldn't be grafted together — they're too different.  Taking that as a challenge, I've tried to adapt the Rx metaphor to Go's concurrency model in a way that does it justice.  For example, `.Subscribe()` returns a channel and a disposable, instead of accepting a callback and returning a disposable:

```go
out, dispose := blah.Subscribe()

x := <-out.Out() // receive from the subscribed channel
dispose.Cancel() // cancel the goroutine managing the subscription
```

# help

As in: please let me know if you'd like to offer some :)  Pull requests encouraged.

# Authors / contributors

- Bryn Bellomy (<bryn@listenonrepeat.com>)