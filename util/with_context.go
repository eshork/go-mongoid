package util

import (
	"context"
	"sync"
	"time"
)

// ContextWithContext returns a new Context that is a functional merging of the two source contexts.
// The resulting context is done whenever either of the sources are done, and Err() will return the value
// as given by the source context that was first discovered as being done.
// Value() will always search within the left context first, then followed by the right context.
func ContextWithContext(left, right context.Context) context.Context {
	c := newMergedCtx(left, right)
	return &c
}

func newMergedCtx(left, right context.Context) mergedCtx {
	return mergedCtx{
		Context: left,
		right:   right,
	}
}

type mergedCtx struct {
	context.Context                 // left context.Context (default)
	right           context.Context // right context.Context
	mu              sync.Mutex      // protects following fields
	first           context.Context // first context.Context to complete
	done            <-chan struct{} // created lazily, closed by first cancel call
}

func (m *mergedCtx) Err() error {
	var err error
	m.mu.Lock()
	if m.first == nil {
		m.mu.Unlock()
		m.Done()
		m.mu.Lock()
	}
	if m.first != nil {
		err = m.first.Err()
		m.mu.Unlock()
		return err
	}
	m.mu.Unlock()
	return err
}

func (m *mergedCtx) Done() <-chan struct{} {
	m.mu.Lock()

	if m.first != nil {
		d := m.first.Done()
		m.mu.Unlock()
		return d
	}

	if m.done != nil {
		d := m.done
		m.mu.Unlock()
		return d
	}

	// log.Fatal("test")

	// if only one side has a Done() that is cancellable, then we can save some effort
	if ldone := m.Context.Done(); ldone != nil {
		if rdone := m.right.Done(); rdone != nil {
			// both have a done channel...

			// see if one is already done, we could just use that
			select {
			case <-ldone: // try left first
				m.first = m.Context
				d := m.first.Done()
				m.mu.Unlock()
				return d
			case <-rdone: // try right next
				m.first = m.right
				d := m.first.Done()
				m.mu.Unlock()
				return d
			default: // default case is breakout if others would block
				// neither channel is done, so we have to monitor them both now
				channel := make(chan struct{})
				m.done = channel
				go func() {
					select {
					case <-ldone:
						m.beDoneVia(m.Context)
						close(channel)
					case <-rdone:
						m.beDoneVia(m.right)
						close(channel)
					}
				}()
				d := m.done
				m.mu.Unlock()
				return d
			}
		} else {
			// only left has a done channel
			m.first = m.Context
			d := m.first.Done()
			m.mu.Unlock()
			return d
		}
	}
	// left isnt cancellable, so whatever right wants to do is fine
	m.first = m.right
	d := m.right.Done()
	m.mu.Unlock()
	return d
}

// makes this mergedCtx done via the result of the given context
func (m *mergedCtx) beDoneVia(c context.Context) {
	m.mu.Lock()
	m.first = c
	m.mu.Unlock()
}

func (m *mergedCtx) Value(key interface{}) interface{} {
	res := m.Context.Value(key)
	if res != nil {
		return res
	}
	return m.right.Value(key)
}

func (m *mergedCtx) Deadline() (deadline time.Time, ok bool) {
	ldl, lok := m.Context.Deadline()
	rdl, rok := m.right.Deadline()
	if lok == true {
		if rok == true {
			if ldl.Before(rdl) {
				return ldl, lok
			}
			return rdl, rok
		}
		return ldl, lok
	}
	return rdl, rok
}
