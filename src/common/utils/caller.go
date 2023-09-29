package utils

import (
	_const "aesir/src/common/const"
	"runtime"
)

func CallerName(skip int) string {
	pcs := make([]uintptr, 1)
	n := runtime.Callers(skip+2, pcs)
	if n < 1 {
		return _const.Unknown
	}

	frame, _ := runtime.CallersFrames(pcs).Next()
	if frame.Function == "" {
		return _const.Unknown
	}
	return frame.Function
}
