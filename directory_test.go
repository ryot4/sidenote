package main

import (
	"testing"
)

func TestJoinPath(t *testing.T) {
	noteDir := ".notes"
	dir := NewDirectory(noteDir)

	tests := []struct {
		name   string
		expect string
	}{
		{"foo", ".notes/foo"},
		{"/foo", ".notes/foo"},
		{"foo.txt", ".notes/foo.txt"},
		{"foo.", ".notes/foo."},
		{"foo/bar", ".notes/foo/bar"},
		{"/foo//bar", ".notes/foo/bar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := dir.JoinPath(tt.name)
			if err != nil {
				t.Error(err)
			}
			if path != tt.expect {
				t.Errorf("expect %q, got %q", tt.expect, path)
			}
		})
	}
}

func TestJoinPathDotFileError(t *testing.T) {
	dir := NewDirectory("/path/to/note")

	tests := []string{
		".foo",
		"/.foo",
		"../foo",
		"foo/.bar",
		"foo/.bar/fizz",
		"foo/./bar",
		"foo/../bar",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			_, err := dir.JoinPath(tt)
			_, ok := err.(*DotFileError)
			if !ok {
				t.Errorf("expect DotFileError, got %v", err)
			}
		})
	}
}
