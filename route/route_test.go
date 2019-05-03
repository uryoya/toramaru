package route

import (
	"testing"
)

func TestRoute_Match(t *testing.T) {
	r := Route{"/a/b", "localhost:8070"}

	matchPaths := []string{"/a/b/c", "/a/basashi"}
	for _, path := range matchPaths {
		if !r.Match(path) {
			t.Errorf("Match(%s) must match %s, got false", path, r.Location)
		}
	}

	unmatchPaths := []string{"/a", "/b/c", "/abc"}
	for _, path := range unmatchPaths {
		if r.Match(path) {
			t.Errorf("Match(%s) must not match %s, got true", path, r.Location)
		}
	}
}
