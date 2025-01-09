package convert

import (
	"os"
	"testing"
)

func TestMp42Flv_Copy(t *testing.T) {
	mp4File, err := os.Open("F:\\test\\test.mp4")
	if err != nil {
		t.Error(err)
		return
	}
	defer mp4File.Close()

	flvFile, err := os.Create("F:\\test\\test.flv")
	if err != nil {
		t.Error(err)
		return
	}
	defer flvFile.Close()

	err = Mp42Flv{}.Copy(flvFile, mp4File)
	if err != nil {
		t.Error(err)
		return
	}
}
