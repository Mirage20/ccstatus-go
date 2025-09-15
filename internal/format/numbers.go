package format

import (
	"fmt"
	"strconv"
)

// WithUnit formats a number with k/M/B suffixes
// Examples:
//
//	1234567 -> "1.2M"
//	45678   -> "46k"
//	789     -> "789"
func WithUnit(value int64) string {
	switch {
	case value >= 1000000000:
		return fmt.Sprintf("%.1fB", float64(value)/1000000000)
	case value >= 1000000:
		return fmt.Sprintf("%.1fM", float64(value)/1000000)
	case value >= 1000:
		return fmt.Sprintf("%dk", value/1000)
	default:
		return strconv.FormatInt(value, 10)
	}
}
