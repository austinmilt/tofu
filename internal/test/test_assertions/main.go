package test_assertions

import (
	"fmt"
	"runtime"
	"testing"
)

type Assert struct {
	t *testing.T
}

func New(t *testing.T) Assert {
	return Assert{t: t}
}

func (a *Assert) True(condition bool, errorMsg string) {
	if !condition {
		a.t.Fatalf("%s expected true but got false: %s", lineRef(), errorMsg)
	}
}

func (a *Assert) False(condition bool, errorMsg string) {
	if condition {
		a.t.Fatalf("%s expected false but got true: %s", lineRef(), errorMsg)
	}
}

func (a *Assert) StringsEqual(expected string, actual string) {
	if expected != actual {
		a.t.Fatalf("%s expected %s but got %s", lineRef(), expected, actual)
	}
}

func (a *Assert) IntsEqual(expected int, actual int) {
	if expected != actual {
		a.t.Fatalf("%s expected %d but got %d", lineRef(), expected, actual)
	}
}

func (a *Assert) IntsOrNilEqual(expected *int, actual *int) {
	if expected != actual {
		a.t.Fatalf("%s expected %d but got %d", lineRef(), expected, actual)
	}
}

func (a *Assert) Uints32Equal(expected uint32, actual uint32) {
	if expected != actual {
		a.t.Fatalf("%s expected %d but got %d", lineRef(), expected, actual)
	}
}

func (a *Assert) Uints32OrNilEqual(expected *uint32, actual *uint32) {
	if expected != actual {
		a.t.Fatalf("%s expected %d but got %d", lineRef(), expected, actual)
	}
}

func (a *Assert) ErrorIsNil(err error) {
	if err != nil {
		a.t.Fatalf("%s expected nil error: %v", lineRef(), err)
	}
}

func (a *Assert) ErrorNotNil(err error) {
	if err == nil {
		a.t.Fatalf("%s expected non-nil error: %v", lineRef(), err)
	}
}

func lineRef() string {
	// use index of 1 to go two levels up the stack
	// (where the assertion would have been called)
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", file, line)
}
