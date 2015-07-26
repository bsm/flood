package health

import (
	"sync/atomic"
	"time"

	"gopkg.in/tomb.v2"
)

// Ping is a continous ping health check
type Ping struct {
	pinger     func() error
	inter      time.Duration
	rise, fall int

	successes, fails, healthy int32

	closer tomb.Tomb
}

// NewPing creates a continous ping health check.
//
// The inter parameter sets the interval between checks.
// The rise parameter sets the number of subsequent checks the ping must pass to be declared healthy.
// The fall parameter sets the number of subsequent failures that would mark the ping as unhealthy.
func NewPing(pinger func() error, inter time.Duration, rise, fall int) *Ping {
	ping := &Ping{
		pinger: pinger,
		inter:  inter,
		rise:   rise,
		fall:   fall,
	}
	ping.closer.Go(ping.loop)
	return ping
}

// IsHealthy implements Check interface
func (p *Ping) IsHealthy() bool {
	return atomic.LoadInt32(&p.healthy) > 0
}

// Stop stops the pinger
func (p *Ping) Stop() {
	p.closer.Kill(nil)
	p.closer.Wait()
}

func (p *Ping) loop() error {
	for {
		select {
		case <-p.closer.Dying():
			return nil
		case <-time.After(p.inter):
			p.update(p.pinger() == nil)
		}
	}
}

func (p *Ping) update(ok bool) {
	if ok {
		atomic.StoreInt32(&p.fails, 0)

		n := int(atomic.AddInt32(&p.successes, 1))
		if n > p.rise {
			atomic.AddInt32(&p.successes, -1)
		} else if n == p.rise {
			atomic.CompareAndSwapInt32(&p.healthy, 0, 1)
		}
	} else {
		atomic.StoreInt32(&p.successes, 0)

		if n := int(atomic.AddInt32(&p.fails, 1)); n > p.fall {
			atomic.AddInt32(&p.fails, -1)
		} else if n == p.fall {
			atomic.CompareAndSwapInt32(&p.healthy, 1, 0)
		}
	}
}
