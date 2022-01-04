package iplink

import "testing"

func Test_runIP(t *testing.T) {
	err := runIP("addr show")
	if err != nil {
		t.Fatal("Failed to run ip command!")
	}

	err = runIP("harp")
	if err != nil {
		t.Log("Tried to run ip command with a wrong argument")
	}
}
