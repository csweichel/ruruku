package kvsession

import (
	"testing"
)

func TestGetLastSegment(t *testing.T) {
	if getLastSegment([]byte("")) != "" {
		t.Errorf("getLastSegment cannot handle empty strings")
	}

	if getLastSegment([]byte(pathSeparator)) != "" {
		t.Errorf("getLastSegment cannot handle malformed paths")
	}

	if getLastSegment(pathSession("foo")) != "foo" {
		t.Errorf("getLastSegment returns the wrong segment")
	}
}
