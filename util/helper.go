package util

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func FormatCurrency(amount float64, currency string) string {
	s := fmt.Sprintf("%.0f", math.Abs(amount))

	var buf strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			buf.WriteByte('.')
		}
		buf.WriteRune(c)
	}

	prefix := ""
	if amount < 0 {
		prefix = "-"
	}
	return prefix + currency + " " + buf.String()
}

func DateFormat(t time.Time) string {
	return t.Format("02 January 2006")
}

func StatusClassCSS(status string) string {
	switch status {
	case "PAID":
		return "status-paid"
	case "PENDING":
		return "status-pending"
	case "EXPIRED":
		return "status-expired"
	default:
		return "status-cancelled"
	}
}
