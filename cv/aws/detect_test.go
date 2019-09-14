package aws

import (
	"os"
	"testing"
)

func TestFaceDetect(t *testing.T) {
	file, err := os.Open("face_test.jpeg")
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
