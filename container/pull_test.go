package container

import "testing"

// TestUnpack is medium test.
func TestUnpack(t *testing.T) {
	c, err := NewClient("./state")
	if err != nil {
		t.Fatalf("NewClinet() returns err=%v", err)
	}
	defer c.Clear()

	if err := c.Pull("alpine"); err != nil {
		t.Errorf("c.Unpack(%s) returns error=%v", "alpine", err)
	}
}
