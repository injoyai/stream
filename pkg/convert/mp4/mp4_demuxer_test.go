package mp4

import (
	"os"
	"testing"
)

func TestNewDemuxer(t *testing.T) {
	mp4File, err := os.Open("F:\\test\\x36xhzz.mp4")
	if err != nil {
		t.Error(err)
		return
	}
	defer mp4File.Close()
	NewDemuxer(mp4File)
}
