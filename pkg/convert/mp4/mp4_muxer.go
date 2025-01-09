package mp4

import (
	"errors"
	"github.com/injoyai/stream/pkg/convert/types"
	"github.com/yapingcat/gomedia/go-mp4"
	"io"
)

func NewMuxer(w io.WriteSeeker) (*Muxer, error) {
	muxer, err := mp4.CreateMp4Muxer(w)
	if err != nil {
		return nil, err
	}
	return &Muxer{
		Movmuxer: muxer,
	}, nil
}

type Muxer struct {
	*mp4.Movmuxer
	hasAudio bool
	hasVideo bool
	atid     uint32
	vtid     uint32
}

func (this *Muxer) WritePacket(p *types.Packet) (err error) {

	switch p.Cid {
	case types.CODEC_H264:
		if !this.hasVideo {
			this.vtid = this.Movmuxer.AddVideoTrack(mp4.MP4_CODEC_H264)
			this.hasVideo = true
		}
		err = this.Movmuxer.Write(this.vtid, p.Data, p.Pts, p.Dts)

	case types.CODEC_AAC:
		if !this.hasAudio {
			this.atid = this.Movmuxer.AddAudioTrack(mp4.MP4_CODEC_AAC)
			this.hasAudio = true
		}
		err = this.Movmuxer.Write(this.atid, p.Data, p.Pts, p.Dts)

	case types.CODEC_MP3:
		if !this.hasAudio {
			this.atid = this.Movmuxer.AddAudioTrack(mp4.MP4_CODEC_MP3)
			this.hasAudio = true
		}
		err = this.Movmuxer.Write(this.atid, p.Data, p.Pts, p.Dts)

	default:
		err = errors.New("unknown codec")

	}

	return
}
