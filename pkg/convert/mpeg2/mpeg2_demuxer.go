package mpeg2

import (
	"github.com/injoyai/base/safe"
	"github.com/injoyai/stream/pkg/convert/types"
	"github.com/yapingcat/gomedia/go-mpeg2"
	"io"
)

func NewDemuxer(r io.ReadSeeker) (*Demuxer, error) {
	demuxer := mpeg2.NewTSDemuxer()
	ch := make(chan *types.Packet)
	demuxer.OnFrame = func(cid mpeg2.TS_STREAM_TYPE, frame []byte, pts uint64, dts uint64) {
		ch <- &types.Packet{
			Cid: func() types.CODEC {
				switch cid {
				case mpeg2.TS_STREAM_AUDIO_MPEG1:
					return types.CODEC_MP3
				case mpeg2.TS_STREAM_AUDIO_MPEG2:
					return types.CODEC_MP3
				case mpeg2.TS_STREAM_AAC:
					return types.CODEC_AAC
				case mpeg2.TS_STREAM_H264:
					return types.CODEC_H264
				case mpeg2.TS_STREAM_H265:
					return types.CODEC_H265
				default:
					return types.CODEC_UNKNOWN
				}
			}(),
			Data:    frame,
			TrackId: 0,
			Pts:     pts,
			Dts:     dts,
		}
	}
	closer := safe.NewCloser()
	go func() {
		err := demuxer.Input(r)
		closer.CloseWithErr(err)
	}()
	return &Demuxer{
		TSDemuxer: demuxer,
		ch:        ch,
		closer:    closer,
	}, nil
}

type Demuxer struct {
	*mpeg2.TSDemuxer
	ch     chan *types.Packet
	closer *safe.Closer
}

func (this *Demuxer) ReadPacket() (*types.Packet, error) {
	select {
	case <-this.closer.Done():
		return nil, this.closer.Err()
	case p := <-this.ch:
		return p, nil
	}
}

func (this *Demuxer) WriteTo(w types.Writer) error {
	for {
		p, err := this.ReadPacket()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err := w.WritePacket(p); err != nil {
			return err
		}
	}
}
