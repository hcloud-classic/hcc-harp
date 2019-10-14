package uuidgen

import "testing"

func Test_Uuidgen(t *testing.T) {
	_, err := UUIDgen()
	if err != nil {
		t.Fatal("Failed to generate uuid!")
	}
}
