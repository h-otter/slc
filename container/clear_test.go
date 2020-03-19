package container

import "testing"

func TestClear(t *testing.T) {
	c, err := NewClient("./state")
	if err != nil {
		t.Fatalf("NewClinet() returns err=%v", err)
	}
	if err :=  c.Clear(); err != nil {
		t.Errorf("Clear() returns err=%v", err)
	}
}
