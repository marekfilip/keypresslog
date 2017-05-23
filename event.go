package keypresslog

import "syscall"

type Event struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

func (e *Event) ToString() string {
	if val, ok := keyCodeMap[e.Code]; ok {
		return val
	}

	return ""
}