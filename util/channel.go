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

// check if a provided channel is closed
func IsChannelClosed[T any](ch <-chan T) bool {
	select {
	case <-ch:
		return false
	default:
		return true
	}
}
