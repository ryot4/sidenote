package main

import (
	"testing"
)

func TestAbsPath(t *testing.T) {
	noteDir := "/note"
	dir := newDirectory(noteDir)

	tests := []struct {
		path   string
		expect string
	}{
		{"foo", "/note/foo"},
		{"/foo", "/note/foo"},
		{"foo.txt", "/note/foo.txt"},
		{"foo.", "/note/foo."},
		{"foo/bar", "/note/foo/bar"},
		{"/foo//bar", "/note/foo/bar"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			path, err := dir.AbsPath(tt.path)
			if err != nil {
				t.Error(err)
			} else if path != tt.expect {
				t.Errorf("expect %q, got %q", tt.expect, path)
			}
		})
	}
}

func TestAbsPathDotFileError(t *testing.T) {
	dir := newDirectory("/path/to/note")

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
			_, err := dir.AbsPath(tt)
			_, ok := err.(*DotFileError)
			if !ok {
				t.Errorf("expect DotFileError, got %v", err)
			}
		})
	}
}
