package types

import "io"

type Demuxer interface {
	Reader
	WriteTo(w Writer) error
}

type Writer interface {
	WritePacket(p *Packet) error
}

type Reader interface {
	ReadPacket() (*Packet, error)
}

type Packet struct {
	Cid     CODEC
	Data    []byte
	TrackId int
	Pts     uint64
	Dts     uint64
}

func Copy(w Writer, r Reader) error {
	for {
		p, err := r.ReadPacket()
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

const (
	CODEC_UNKNOWN CODEC = iota
	CODEC_H264
	CODEC_H265
	CODEC_VP8
)

const (
	CODEC_AAC CODEC = iota + 101
	CODEC_MP2
	CODEC_MP3
	CODEC_G711A
	CODEC_G711U
	CODEC_OPUS
)

type CODEC int

func (c CODEC) String() string {
	switch c {
	case CODEC_H264:
		return "h264"
	case CODEC_H265:
		return "h265"
	case CODEC_AAC:
		return "aac"
	case CODEC_MP2:
		return "mp2"
	case CODEC_MP3:
		return "mp3"
	case CODEC_G711A:
		return "g711a"
	case CODEC_G711U:
		return "g711u"
	case CODEC_OPUS:
		return "opus"
	default:
		return "unknown"
	}
}
