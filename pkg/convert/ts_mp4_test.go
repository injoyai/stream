package convert

import (
	"os"
	"testing"
)

func TestTs2Mp4_Copy(t *testing.T) {
	tsFile, err := os.Open("F:\\test\\test.ts")
	if err != nil {
		t.Error(err)
		return
	}
	defer tsFile.Close()

	mp4File, err := os.Create("F:\\test\\test2.mp4")
	if err != nil {
		t.Error(err)
		return
	}
	defer mp4File.Close()

	err = Ts2Mp4{}.Copy(mp4File, tsFile)
	if err != nil {
		t.Error(err)
		return
	}
}
