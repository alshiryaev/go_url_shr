package link

import (
	"testing"
)

func TestRandomStringRunes(t *testing.T) {
	got := RandomStringRunes(10)
	if len(got) != 10 {
		t.Errorf("wanted 10 but got=%d", len(got))
	}
}
