package mpeg2

import (
	"errors"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/stream/pkg/convert/types"
	"github.com/yapingcat/gomedia/go-mpeg2"
	"io"
)

func NewMuxer(w io.Writer) *Muxer {
	muxer := mpeg2.NewTSMuxer()
	closer := safe.NewCloser()
	muxer.OnPacket = func(pkg []byte) {
		_, err := w.Write(pkg)
		closer.CloseWithErr(err)
	}
	return &Muxer{
		TSMuxer: muxer,
		closer:  closer,
	}
}

type Muxer struct {
	*mpeg2.TSMuxer
	hasAudio bool
	hasVideo bool
	atid     uint16
	vtid     uint16
	closer   *safe.Closer
}

func (this *Muxer) WritePacket(p *types.Packet) error {
	switch p.Cid {
	case types.CODEC_H264:
		if !this.hasVideo {
			this.vtid = this.TSMuxer.AddStream(mpeg2.TS_STREAM_H264)
			this.hasVideo = true
		}
		if err := this.TSMuxer.Write(this.vtid, p.Data, p.Pts, p.Dts); err != nil {
			return err
		}

	case types.CODEC_H265:
		if !this.hasVideo {
			this.vtid = this.TSMuxer.AddStream(mpeg2.TS_STREAM_H265)
			this.hasVideo = true
		}
		if err := this.TSMuxer.Write(this.vtid, p.Data, p.Pts, p.Dts); err != nil {
			return err
		}

	case types.CODEC_MP3:
		if !this.hasAudio {
			this.hasAudio = true
			this.atid = this.TSMuxer.AddStream(mpeg2.TS_STREAM_AUDIO_MPEG2)
		}
		if err := this.TSMuxer.Write(this.atid, p.Data, p.Pts, p.Dts); err != nil {
			return err
		}

	case types.CODEC_AAC:
		if !this.hasAudio {
			this.hasAudio = true
			this.atid = this.TSMuxer.AddStream(mpeg2.TS_STREAM_AAC)
		}
		if err := this.TSMuxer.Write(this.atid, p.Data, p.Pts, p.Dts); err != nil {
			return err
		}

	default:
		return errors.New("unknown codec: " + p.Cid.String())

	}
	return this.closer.Err()
}
