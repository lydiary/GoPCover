package views

import "testing"

func TestGetRevisionFromVersion(t *testing.T) {
	version := "1.1.1.1"
	expect := "1234"
	got := GetRevisionFromVersion(version)
	if expect != got {
		t.Errorf("expect %s, got %s", expect, got)
	}
}
