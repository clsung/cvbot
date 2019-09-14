package gcp

import (
	"os"
	"testing"
)

func TestFaceDetect(t *testing.T) {
	file, err := os.Open("1559725685111.jpg")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer file.Close()
	retStr, err := FaceDetect(file)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(retStr) < 10 {
		t.Errorf("%s", retStr)
	}
}
