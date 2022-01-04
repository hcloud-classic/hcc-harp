package cmd

import "testing"

func Test_runCMD(t *testing.T) {
	err := RunCMD("ls -al")
	if err != nil {
		t.Fatal("Failed to run ls command!")
	}

	err = RunCMD("9 8 7 6 5 4 3 2 1")
	if err != nil {
		t.Log("Tried to run wrong command")
	}
}
