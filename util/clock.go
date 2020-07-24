package util

import "time"

// Clock abstracts the system clock
type Clock interface {
	CurrentTime() time.Time
}

// SystemClock clock returns the time in local time zone
type SystemClock struct{}

// CurrentTime returns the current time in local time zone
func (SystemClock) CurrentTime() time.Time {
	return time.Now()
}
