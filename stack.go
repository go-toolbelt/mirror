package mirror

import (
	"runtime"
)

const depth = 32

type Stack struct {
	ptrs [depth]uintptr
	n    int
}

// baseSkips is an additional number of skips that must be passed to get the right
// traces. It was derived empirically through testing.
const baseSkips = 1

//go:noinline
func Capture(skip int) Stack {
	var stack Stack
	stack.n = runtime.Callers(skip+baseSkips+1, stack.ptrs[:])
	return stack
}

type Frames struct {
	ptrs    []uintptr
	current []Frame
}

func (stack *Stack) Frames() Frames {
	return Frames{
		ptrs: stack.ptrs[:stack.n],
	}
}

func (frames *Frames) Next() (frame Frame, ok bool) {
	for len(frames.current) == 0 {
		if len(frames.ptrs) == 0 {
			return Frame{}, false
		}

		frames.current, frames.ptrs = getFrameForPtr(frames.ptrs[0]), frames.ptrs[1:]
	}

	frame, frames.current = frames.current[0], frames.current[1:]
	return frame, true
}
