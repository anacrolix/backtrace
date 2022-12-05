package backtrace

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
)

// Adds a backtrace to the error if an existing one isn't found. It's possible that multiple
// backtraces might be added in the future.
func With(err error) error {
	if err == nil {
		return nil
	}
	var with errorWith
	if errors.As(err, &with) {
		return err
	}
	with.numUsed = runtime.Callers(2, with.callers[:])
	with.wrapped = err
	return with
}

type errorWith struct {
	callers [32]uintptr
	numUsed int
	wrapped error
}

func (me errorWith) Error() string {
	return me.wrapped.Error()
}

func (me errorWith) Unwrap() error {
	return me.wrapped
}

func (me errorWith) toBacktrace() Backtrace {
	return Backtrace{
		Callers: me.callers[:me.numUsed],
		Ok:      true,
	}
}

// Retrieves the first backtrace found. I suspect multiple backtraces might exist at some point.
// It's reserved that this might change to return the most appropriate level of backtrace.
func FromError(err error) (bt Backtrace) {
	var with errorWith
	if !errors.As(err, &with) {
		return
	}
	return with.toBacktrace()
}

// A backtrace that can be moved around inside a process, such as included with errors for later
// extraction.
type Backtrace struct {
	// Contains the callers for the location where this backtrace was added.
	Callers []uintptr
	// False if no backtrace exists.
	Ok bool
}

// Prepends a newline if Ok, otherwise returns an empty string. Intended to be inserted into log
// formatting fields so the backtrace is added on its own lines if appropriate.
func (me Backtrace) TryInject() string {
	if !me.Ok {
		return ""
	}
	var builder strings.Builder
	builder.WriteRune('\n')
	me.build(&builder)
	return builder.String()
}

// I wonder if this should a concise, single-line form.
func (me Backtrace) String() string {
	var builder strings.Builder
	me.build(&builder)
	return builder.String()
}

func (me Backtrace) build(builder *strings.Builder) {
	builder.WriteString("  backtrace:\n")
	frames := runtime.CallersFrames(me.Callers)
	for {
		frame, more := frames.Next()
		builder.WriteString("    ")
		builder.WriteString(frame.Function)
		builder.WriteRune(':')
		builder.WriteString(strconv.Itoa(frame.Line))
		builder.WriteRune('\n')
		if !more {
			break
		}
	}
}

// Outputs an empty string if there is no backtrace, otherwise outputs the backtrace prepended with
// a newline for injecting into other messages.
func Sprint(err error) string {
	return FromError(err).TryInject()
}

// Should we have a helper that does backtrace and wrapping in one?
//func Errorf()
