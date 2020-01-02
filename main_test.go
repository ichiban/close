package close

import (
	"errors"
	"io"
	"testing"
)

func TestLater(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		closers = nil
		var c mock
		assertNotPanics(t, func() {
			Later(&c)
		})
		assertContains(t, closers, &c)
	})

	t.Run("single", func(t *testing.T) {
		var c1, c2 mock
		closers = []io.Closer{&c1}
		assertNotPanics(t, func() {
			Later(&c2)
		})
		assertContains(t, closers, &c1)
		assertContains(t, closers, &c2)
	})
}

func TestAll(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		closers = nil
		assertNotPanics(t, func() {
			All()
		})
	})

	t.Run("succeed", func(t *testing.T) {
		var c mock
		closers = []io.Closer{&c}
		assertNotPanics(t, func() {
			All()
		})
		assertTrue(t, c.closed)
	})

	t.Run("fail", func(t *testing.T) {
		var c1, c3 mock
		c2 := mock{
			err: errors.New(""),
		}
		closers = []io.Closer{&c1, &c2, &c3}
		assertPanics(t, func() {
			All()
		})
		assertTrue(t, c1.closed)
		assertTrue(t, c2.closed)
		assertTrue(t, c3.closed)
	})
}

type mock struct {
	closed bool
	err    error
}

func (m *mock) Close() error {
	m.closed = true
	return m.err
}

func assertNotPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should not panic: %v", r)
		}
	}()
	f()
}

func assertPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("should panic: %v", r)
		}
	}()
	f()
}

func assertContains(t *testing.T, s []io.Closer, contains io.Closer) {
	t.Helper()
	for _, e := range s {
		if e == contains {
			return
		}
	}
	t.Errorf("should contain: %v ∋ %v", s, contains)
}

func assertTrue(t *testing.T, o bool) {
	t.Helper()
	if !o {
		t.Errorf("should be true: %v", o)
	}
}
