package util

// check if a provided channel is closed
func IsChannelClosed[T any](ch <-chan T) bool {
	select {
	case <-ch:
		return false
	default:
		return true
	}
}
