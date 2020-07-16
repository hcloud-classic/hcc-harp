package fileutil

import "testing"

func Test_WriteFile(t *testing.T) {
	err := WriteFile("/tmp/harp_test_file", "test")
	if err != nil {
		t.Fatal("Failed to write file!")
	}
}

func Test_WriteFileAppend(t *testing.T) {
	err := WriteFileAppend("/tmp/harp_test_file", "_test_append")
	if err != nil {
		t.Fatal("Failed to write file in append mode!")
	}
}

func Test_DeleteFile(t *testing.T) {
	err := DeleteFile("/tmp/harp_test_file")
	if err != nil {
		t.Error("Failed to delete file!")
	}
}