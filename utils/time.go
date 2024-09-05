package utils

import "time"

// IsDawn 是否是凌晨
func IsDawn(t time.Time) bool {
	return t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0
}
