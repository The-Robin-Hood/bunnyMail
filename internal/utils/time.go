package utils

import "time"

func ToISO(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
