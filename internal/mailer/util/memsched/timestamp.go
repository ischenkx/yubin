package memsched

import "time"

type TimeStamp[T any] struct {
	ID   string
	Data T
	Time time.Time
}
