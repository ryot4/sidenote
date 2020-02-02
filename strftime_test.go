package main

import (
	"testing"
	"time"
)

func TestStrftime(t *testing.T) {
	refTime, err := time.Parse(time.UnixDate, "Sat Feb  3 16:05:06 JST 2001")
	if err != nil {
		panic(err)
	}

	tests := []struct {
		format string
		expect string
	}{
		{"%a", "Sat"},
		{"%A", "Saturday"},
		{"%b", "Feb"},
		{"%B", "February"},
		{"%d", "03"},
		{"%e", " 3"},
		{"%F", "2001-02-03"},
		{"%H", "16"},
		{"%I", "04"},
		{"%m", "02"},
		{"%M", "05"},
		{"%p", "PM"},
		{"%S", "06"},
		{"%y", "01"},
		{"%Y", "2001"},
		{"%z", "+0900"},
		{"%Z", "JST"},
		{"%%", "%"},

		{"%a %b %e %H:%M:%S %Z %Y", "Sat Feb  3 16:05:06 JST 2001"},
		{"diary/%Y/%m/%d.txt", "diary/2001/02/03.txt"},
		{"%Y年%m月%d日_%H時%M分%S秒.md", "2001年02月03日_16時05分06秒.md"},
		{"test.txt", "test.txt"},
		{"%TODO", "%TODO"},
		{"TODO%", "TODO%"},
		{"%%TODO%%", "%TODO%"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			formatted := Strftime(refTime, tt.format)
			if formatted != tt.expect {
				t.Errorf("expect %q, got %q", tt.expect, formatted)
			}
		})
	}
}
