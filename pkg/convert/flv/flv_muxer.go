package flv

import (
	"fmt"
	"github.com/injoyai/stream/pkg/convert/types"
	"github.com/yapingcat/gomedia/go-flv"
	"io"
)

func NewMuxer(w io.Writer) (*Muxer, error) {
	muxer := flv.CreateFlvWriter(w)
	if err := muxer.WriteFlvHeader(); err != nil {
		return nil, err
	}
	return &Muxer{
		FlvWriter: muxer,
	}, nil
}

type Muxer struct {
	*flv.FlvWriter
}

func (this *Muxer) WritePacket(p *types.Packet) (err error) {
	switch p.Cid {
	case types.CODEC_H264:
		err = this.FlvWriter.WriteH264(p.Data, uint32(p.Pts), uint32(p.Dts))

	case types.CODEC_H265:
		err = this.FlvWriter.WriteH265(p.Data, uint32(p.Pts), uint32(p.Dts))

	case types.CODEC_AAC:
		err = this.FlvWriter.WriteAAC(p.Data, uint32(p.Pts), uint32(p.Dts))

	case types.CODEC_G711A:
		err = this.FlvWriter.WriteG711A(p.Data, uint32(p.Pts), uint32(p.Dts))

	case types.CODEC_G711U:
		err = this.FlvWriter.WriteG711U(p.Data, uint32(p.Pts), uint32(p.Dts))

	case types.CODEC_MP3:
		err = this.FlvWriter.WriteMp3(p.Data, uint32(p.Pts), uint32(p.Dts))

	default:
		err = fmt.Errorf("未知编码格式: %X", p.Cid)

	}
	return
}
