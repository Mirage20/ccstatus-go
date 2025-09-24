package format

import (
	"fmt"
	"strconv"
)

const (
	// Units for number formatting.
	million = 1000000
	billion = 1000000000
)

// WithUnit formats a number with k/M/B suffixes
// Examples:
//
//	1234567 -> "1.2M"
//	45678   -> "46k"
//	789     -> "789"
func WithUnit(value int64) string {
	switch {
	case value >= billion:
		return fmt.Sprintf("%.1fB", float64(value)/billion)
	case value >= million:
		return fmt.Sprintf("%.1fM", float64(value)/million)
	case value >= 1000:
		return fmt.Sprintf("%dk", value/1000)
	default:
		return strconv.FormatInt(value, 10)
	}
}
