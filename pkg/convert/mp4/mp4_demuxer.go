package mp4

import (
	"github.com/injoyai/stream/pkg/convert/types"
	"github.com/yapingcat/gomedia/go-mp4"
	"io"
)

func NewDemuxer(r io.ReadSeeker) (*Demuxer, error) {
	d := mp4.CreateMp4Demuxer(r)
	tracks, err := d.ReadHead()
	if err != nil {
		return nil, err
	}
	return &Demuxer{
		MovDemuxer: d,
		Tracks:     tracks,
	}, nil
}

type Demuxer struct {
	*mp4.MovDemuxer
	Tracks []mp4.TrackInfo //轨道,包括音频和视频
}

func (this *Demuxer) ReadPacket() (*types.Packet, error) {
	p, err := this.MovDemuxer.ReadPacket()
	if err != nil {
		return nil, err
	}
	return &types.Packet{
		Cid: func() types.CODEC {
			switch p.Cid {
			case mp4.MP4_CODEC_H264:
				return types.CODEC_H264
			case mp4.MP4_CODEC_H265:
				return types.CODEC_H265
			case mp4.MP4_CODEC_AAC:
				return types.CODEC_AAC
			case mp4.MP4_CODEC_G711A:
				return types.CODEC_G711A
			case mp4.MP4_CODEC_G711U:
				return types.CODEC_G711U
			case mp4.MP4_CODEC_MP2:
				return types.CODEC_MP2
			case mp4.MP4_CODEC_MP3:
				return types.CODEC_MP3
			case mp4.MP4_CODEC_OPUS:
				return types.CODEC_OPUS
			default:
				return types.CODEC_UNKNOWN
			}
		}(),
		Data:    p.Data,
		TrackId: p.TrackId,
		Pts:     p.Pts,
		Dts:     p.Dts,
	}, nil
}

func (this *Demuxer) WriteTo(w types.Writer) error {
	return types.Copy(w, this)
}
