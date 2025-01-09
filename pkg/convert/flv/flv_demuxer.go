package flv

import (
	"github.com/injoyai/base/safe"
	"github.com/injoyai/stream/pkg/convert/types"
	"github.com/yapingcat/gomedia/go-codec"
	"github.com/yapingcat/gomedia/go-flv"
	"io"
)

func NewDemuxer(r io.Reader) *Demuxer {
	reader := flv.CreateFlvReader()
	ch := make(chan *types.Packet)
	reader.OnFrame = func(cid codec.CodecID, frame []byte, pts uint32, dts uint32) {
		ch <- &types.Packet{
			Cid: func() types.CODEC {
				switch cid {
				case codec.CODECID_VIDEO_H264:
					return types.CODEC_H264
				case codec.CODECID_VIDEO_H265:
					return types.CODEC_H265
				case codec.CODECID_VIDEO_VP8:
					return types.CODEC_VP8
				case codec.CODECID_AUDIO_AAC:
					return types.CODEC_AAC
				case codec.CODECID_AUDIO_G711A:
					return types.CODEC_G711A
				case codec.CODECID_AUDIO_G711U:
					return types.CODEC_G711U
				case codec.CODECID_AUDIO_OPUS:
					return types.CODEC_OPUS
				case codec.CODECID_AUDIO_MP3:
					return types.CODEC_MP3
				default:
					return types.CODEC_UNKNOWN
				}

			}(),
			Data:    frame,
			TrackId: 0,
			Pts:     uint64(pts),
			Dts:     uint64(dts),
		}
	}

	demuxer := &Demuxer{
		FlvReader: reader,
		ch:        ch,
		closer:    safe.NewCloser(),
	}

	go func() {
		_, err := io.Copy(demuxer, r)
		demuxer.closer.CloseWithErr(err)
	}()
	return demuxer
}

type Demuxer struct {
	*flv.FlvReader
	ch     chan *types.Packet
	closer *safe.Closer
}

func (this *Demuxer) Write(p []byte) (int, error) {
	err := this.FlvReader.Input(p)
	return len(p), err
}

func (this *Demuxer) ReadPacket() (*types.Packet, error) {
	select {
	case <-this.closer.Done():
		return nil, this.closer.Err()
	case p := <-this.ch:
		return p, nil
	}
}
