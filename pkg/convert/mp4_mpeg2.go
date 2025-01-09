package convert

import (
	"github.com/yapingcat/gomedia/go-mp4"
	"github.com/yapingcat/gomedia/go-mpeg2"
	"io"
)

type Mp4ToMpeg2 struct{}

func (m Mp4ToMpeg2) Copy(w io.Writer, r io.ReadSeeker) error {

	muxer := mpeg2.NewTSMuxer()
	muxer.OnPacket = func(pkg []byte) { w.Write(pkg) }

	demuxer := mp4.CreateMp4Demuxer(r)
	if _, err := demuxer.ReadHead(); err != nil {
		return err
	}

	hasAudio := false
	hasVideo := false
	var atid uint16 = 0
	var vtid uint16 = 0
	for {
		p, err := demuxer.ReadPacket()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch p.Cid {
		case mp4.MP4_CODEC_H264:
			if !hasVideo {
				vtid = muxer.AddStream(mpeg2.TS_STREAM_H264)
				hasVideo = true
			}
			if err := muxer.Write(vtid, p.Data, p.Pts, p.Dts); err != nil {
				return err
			}

		case mp4.MP4_CODEC_H265:
			if !hasVideo {
				vtid = muxer.AddStream(mpeg2.TS_STREAM_H265)
				hasVideo = true
			}
			if err := muxer.Write(vtid, p.Data, p.Pts, p.Dts); err != nil {
				return err
			}

		case mp4.MP4_CODEC_MP3:
			if !hasAudio {
				hasAudio = true
				atid = muxer.AddStream(mpeg2.TS_STREAM_AUDIO_MPEG2)
			}
			if err := muxer.Write(atid, p.Data, p.Pts, p.Dts); err != nil {
				return err
			}

		case mp4.MP4_CODEC_AAC:
			if !hasAudio {
				hasAudio = true
				atid = muxer.AddStream(mpeg2.TS_STREAM_AAC)
			}
			if err := muxer.Write(atid, p.Data, p.Pts, p.Dts); err != nil {
				return err
			}

		}

	}

}
