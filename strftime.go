package main

import (
	"strings"
	"time"
)

// Strftime returns string representation of t formatted according to format.
func Strftime(t time.Time, format string) string {
	b := strings.Builder{}
	runes := []rune(format)
	end := len(runes)
	b.Grow(end)
	for i := 0; i < end; i++ {
		c := runes[i]
		if c != '%' {
			b.WriteRune(c)
			continue
		}
		i++
		if i >= end {
			// format ends with '%'
			b.WriteRune('%')
			break
		}
		switch runes[i] {
		case 'a':
			b.WriteString(t.Format("Mon"))
		case 'A':
			b.WriteString(t.Format("Monday"))
		case 'b':
			b.WriteString(t.Format("Jan"))
		case 'B':
			b.WriteString(t.Format("January"))
		case 'd':
			b.WriteString(t.Format("02"))
		case 'e':
			b.WriteString(t.Format("_2"))
		case 'F':
			b.WriteString(t.Format("2006-01-02"))
		case 'H':
			b.WriteString(t.Format("15"))
		case 'I':
			b.WriteString(t.Format("03"))
		case 'm':
			b.WriteString(t.Format("01"))
		case 'M':
			b.WriteString(t.Format("04"))
		case 'p':
			b.WriteString(t.Format("PM"))
		case 'S':
			b.WriteString(t.Format("05"))
		case 'y':
			b.WriteString(t.Format("06"))
		case 'Y':
			b.WriteString(t.Format("2006"))
		case 'z':
			b.WriteString(t.Format("-0700"))
		case 'Z':
			b.WriteString(t.Format("MST"))
		case '%':
			b.WriteRune('%')
		default:
			b.WriteRune('%')
			b.WriteRune(runes[i])
		}
	}
	return b.String()
}
