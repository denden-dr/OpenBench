package apierrors

import (
	"fmt"
	"runtime"
)

type withStack struct {
	err   error
	msg   string
	stack []uintptr
}

func (w *withStack) Error() string {
	if w.msg != "" {
		if w.err != nil {
			return w.msg + ": " + w.err.Error()
		}
		return w.msg
	}
	if w.err != nil {
		return w.err.Error()
	}
	return ""
}

func (w *withStack) Unwrap() error {
	return w.err
}

// StackTrace returns the formatted stack trace.
func (w *withStack) StackTrace() string {
	frames := runtime.CallersFrames(w.stack)
	var result string
	for {
		frame, more := frames.Next()
		result += fmt.Sprintf("\n%s\n\t%s:%d", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return result
}

// Wrap wraps an existing error with a message and captures the stack trace.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err:   err,
		msg:   message,
		stack: callers(),
	}
}

// New creates a new error with a message and captures the stack trace.
func New(message string) error {
	return &withStack{
		msg:   message,
		stack: callers(),
	}
}

func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[:n]
}
