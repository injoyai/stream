package convert

import (
	"fmt"
	"github.com/yapingcat/gomedia/go-flv"
	"github.com/yapingcat/gomedia/go-mp4"
	"io"
)

type Mp42Flv struct {
}

func (this Mp42Flv) Copy(w io.Writer, r io.ReadSeeker) error {
	demuxer := mp4.CreateMp4Demuxer(r)
	if _, err := demuxer.ReadHead(); err != nil {
		return err
	}

	writer := flv.CreateFlvWriter(w)
	if err := writer.WriteFlvHeader(); err != nil {
		return err
	}

	for {
		pkg, err := demuxer.ReadPacket()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch pkg.Cid {
		case mp4.MP4_CODEC_H264:
			err = writer.WriteH264(pkg.Data, uint32(pkg.Pts), uint32(pkg.Dts))

		case mp4.MP4_CODEC_H265:
			err = writer.WriteH265(pkg.Data, uint32(pkg.Pts), uint32(pkg.Dts))

		case mp4.MP4_CODEC_AAC:
			err = writer.WriteAAC(pkg.Data, uint32(pkg.Pts), uint32(pkg.Dts))

		case mp4.MP4_CODEC_G711A:
			err = writer.WriteG711A(pkg.Data, uint32(pkg.Pts), uint32(pkg.Dts))

		case mp4.MP4_CODEC_G711U:
			err = writer.WriteG711U(pkg.Data, uint32(pkg.Pts), uint32(pkg.Dts))

		case mp4.MP4_CODEC_MP3:
			err = writer.WriteMp3(pkg.Data, uint32(pkg.Pts), uint32(pkg.Dts))

		default:
			err = fmt.Errorf("未知编码格式: %X", pkg.Cid)

		}

		if err != nil {
			return err
		}

	}
}
