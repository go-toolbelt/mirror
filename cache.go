package mirror

import (
	"runtime"
	"sync"
)

// This global frame cache makes sense because it's store a mapping of memory address
// to stack frame identifier which is guaranteed to never change globally.
// nolint: gochecknoglobals
var frameCache sync.Map // map[uintptr][]Frame

func getFrameForPtr(ptr uintptr) []Frame {
	value, ok := frameCache.Load(ptr)
	if ok {
		return value.([]Frame)
	}

	callers := runtime.CallersFrames([]uintptr{ptr})

	var frames []Frame
	for {
		frame, more := callers.Next()

		frames = append(frames, fromRuntimeFrame(frame))

		if !more {
			break
		}
	}
	frameCache.Store(ptr, frames)
	return frames
}
