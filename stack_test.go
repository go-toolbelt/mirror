package mirror_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/go-toolbelt/benchmark"

	"github.com/go-toolbelt/mirror"
)

const (
	packageName = "github.com/go-toolbelt/mirror_test"
	fileName    = "stack_test.go"
)

func TestCapture(t *testing.T) {
	result := captureStackAndLine()
	frames := result.stack.Frames()
	frame, more := frames.Next()
	assertTrue(t, more)
	assertEqual(t, mirror.Frame{
		Package:   packageName,
		Function:  "TestCapture",
		File:      fileName,
		Line:      result.line,
		Formatted: fmt.Sprintf("\tat %s.TestCapture(%s:%d)", packageName, fileName, result.line),
	}, frame)
}

type stackAndLine struct {
	stack mirror.Stack
	line  int
}

func exampleFunc() stackAndLine {
	return captureStackAndLine()
}

func TestCapture_Func(t *testing.T) {
	// N.B. Both function calls must be kept on the same line to ensure line numbers will be equal
	result, line := exampleFunc(), captureLineNumber(0)

	frames := result.stack.Frames()

	frame, more := frames.Next()
	assertTrue(t, more)
	assertEqual(t, mirror.Frame{
		Package:   packageName,
		Function:  "exampleFunc",
		File:      fileName,
		Line:      result.line,
		Formatted: fmt.Sprintf("\tat %s.exampleFunc(%s:%d)", packageName, fileName, result.line),
	}, frame)
	assertTrue(t, more)

	frame, more = frames.Next()
	assertTrue(t, more)
	assertEqual(t, mirror.Frame{
		Package:   packageName,
		Function:  "TestCapture_Func",
		File:      fileName,
		Line:      line,
		Formatted: fmt.Sprintf("\tat %s.TestCapture_Func(%s:%d)", packageName, fileName, line),
	}, frame)
}

func TestCapture_InnerFunc(t *testing.T) {
	// purposefully testing a lambda here.
	// nolint: gocritic
	exampleInnerFunc := func() stackAndLine {
		return captureStackAndLine()
	}

	// N.B. Both function calls must be kept on the same line to ensure line numbers will be equal
	result, line := exampleInnerFunc(), captureLineNumber(0)

	frames := result.stack.Frames()

	frame, more := frames.Next()
	assertTrue(t, more)
	assertEqual(t, mirror.Frame{
		Package:   packageName,
		Function:  "TestCapture_InnerFunc.func1",
		File:      fileName,
		Line:      result.line,
		Formatted: fmt.Sprintf("\tat %s.TestCapture_InnerFunc.func1(%s:%d)", packageName, fileName, result.line),
	}, frame)
	assertTrue(t, more)

	frame, more = frames.Next()
	assertTrue(t, more)
	assertEqual(t, mirror.Frame{
		Package:   packageName,
		Function:  "TestCapture_InnerFunc",
		File:      fileName,
		Line:      line,
		Formatted: fmt.Sprintf("\tat %s.TestCapture_InnerFunc(%s:%d)", packageName, fileName, line),
	}, frame)
}

type exampleStruct struct{}

func (exampleStruct) exampleStructMethod() stackAndLine {
	return captureStackAndLine()
}

func TestCapture_StructMethod(t *testing.T) {
	// N.B. Both function calls must be kept on the same line to ensure line numbers will be equal
	result, line := exampleStruct{}.exampleStructMethod(), captureLineNumber(0)

	frames := result.stack.Frames()

	frame, more := frames.Next()
	assertTrue(t, more)
	assertEqual(t, mirror.Frame{
		Package:   packageName,
		Function:  "exampleStruct.exampleStructMethod",
		File:      fileName,
		Line:      result.line,
		Formatted: fmt.Sprintf("\tat %s.exampleStruct.exampleStructMethod(%s:%d)", packageName, fileName, result.line),
	}, frame)
	assertTrue(t, more)

	frame, more = frames.Next()
	assertTrue(t, more)
	assertEqual(t, mirror.Frame{
		Package:   packageName,
		Function:  "TestCapture_StructMethod",
		File:      fileName,
		Line:      line,
		Formatted: fmt.Sprintf("\tat %s.TestCapture_StructMethod(%s:%d)", packageName, fileName, line),
	}, frame)
}

type exampleStructPtr struct{}

func (*exampleStructPtr) exampleStructPtrMethod() stackAndLine {
	return captureStackAndLine()
}

func TestCapture_StructPtrMethod(t *testing.T) {
	example := &exampleStructPtr{}

	// N.B. Both function calls must be kept on the same line to ensure line numbers will be equal
	result, line := example.exampleStructPtrMethod(), captureLineNumber(0)

	frames := result.stack.Frames()

	frame, more := frames.Next()
	assertTrue(t, more)
	assertEqual(t, mirror.Frame{
		Package:   packageName,
		Function:  "exampleStructPtr.exampleStructPtrMethod",
		File:      fileName,
		Line:      result.line,
		Formatted: fmt.Sprintf("\tat %s.exampleStructPtr.exampleStructPtrMethod(%s:%d)", packageName, fileName, result.line),
	}, frame)

	frame, more = frames.Next()
	assertTrue(t, more)
	assertEqual(t, mirror.Frame{
		Package:   packageName,
		Function:  "TestCapture_StructPtrMethod",
		File:      fileName,
		Line:      line,
		Formatted: fmt.Sprintf("\tat %s.TestCapture_StructPtrMethod(%s:%d)", packageName, fileName, line),
	}, frame)
}

func TestBenchmark_Frames(t *testing.T) {
	benchmark.Test(
		t,
		func(b *testing.B) {
			// Capture at the same line number for all callers.
			capture := func() mirror.Stack {
				return mirror.Capture(0)
			}

			// Capture the frame once to save it in the cache
			stack := capture()
			frames := stack.Frames()
			expected, _ := frames.Next()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				stack := capture()
				frames := stack.Frames()
				actual, _ := frames.Next()
				if expected != actual {
					t.FailNow()
				}
			}
		},
		benchmark.ZeroAllocsPerOp(),
	)
}

func assertTrue(t *testing.T, b bool) {
	if !b {
		t.Error("Must be true")
	}
}

func assertEqual(t *testing.T, expected mirror.Frame, actual mirror.Frame) {
	if expected != actual {
		t.Errorf("Frames not equal: expected %s actual %s", expected, actual)
	}
}

func captureStackAndLine() stackAndLine {
	return stackAndLine{
		stack: mirror.Capture(1),
		line:  captureLineNumber(1),
	}
}

func captureLineNumber(skip int) int {
	_, _, line, _ := runtime.Caller(skip + 1)
	return line
}
