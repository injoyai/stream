package convert

import (
	"github.com/injoyai/stream/pkg/convert/mp4"
	"github.com/injoyai/stream/pkg/convert/mpeg2"
	"os"
	"testing"
)

func TestMp4Mpeg2_Copy(t *testing.T) {
	tsFile, err := os.Create("F:\\test\\x36xhzz_2.ts")
	if err != nil {
		t.Error(err)
		return
	}
	defer tsFile.Close()

	mp4File, err := os.Open("F:\\test\\x36xhzz.mp4")
	if err != nil {
		t.Error(err)
		return
	}
	defer mp4File.Close()

	err = Mp4ToMpeg2{}.Copy(tsFile, mp4File)
	t.Log(err)
}

func TestMp4_Mpeg2(t *testing.T) {
	tsFile, err := os.Create("F:\\test\\x36xhzz_3.ts")
	if err != nil {
		t.Error(err)
		return
	}
	defer tsFile.Close()

	mp4File, err := os.Open("F:\\test\\x36xhzz.mp4")
	if err != nil {
		t.Error(err)
		return
	}
	defer mp4File.Close()

	demuxer, err := mp4.NewDemuxer(mp4File)
	if err != nil {
		t.Error(err)
		return
	}

	muxer := mpeg2.NewMuxer(tsFile)

	err = demuxer.WriteTo(muxer)

	//err = types.Copy(muxer, demuxer)

	t.Log(err)

}
