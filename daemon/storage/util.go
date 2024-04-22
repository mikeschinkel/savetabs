package storage

import (
	"time"
)

func throttle() {
	time.Sleep(250 * time.Millisecond)
}
