package helpers

import "time"

func ComposeDurations(ttlStr, timeoutStr string) (ttl time.Duration, timeout time.Duration, err error) {
	ttl, err = time.ParseDuration(ttlStr)
	if err != nil {
		return
	}

	timeout, err = time.ParseDuration(timeoutStr)
	return
}
