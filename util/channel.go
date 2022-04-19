package util

import "time"

func ReceiveWithTimeout[T any](c chan T, maxDelay time.Duration) T {
	timeout := time.After(maxDelay)

	select {
	case data := <-c:
		return data
	case <-timeout:
		return Zero[T]()
	}
}
