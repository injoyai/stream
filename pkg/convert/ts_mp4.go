package convert

import (
	"github.com/yapingcat/gomedia/go-mp4"
	"github.com/yapingcat/gomedia/go-mpeg2"
	"io"
)

type Ts2Mp4 struct{}

func (Ts2Mp4) Copy(w io.WriteSeeker, r io.ReadSeeker) error {

	muxer, err := mp4.CreateMp4Muxer(w)
	if err != nil {
		return err
	}

	hasAudio := false
	hasVideo := false
	var atid uint32 = 0
	var vtid uint32 = 0
	demuxer := mpeg2.NewTSDemuxer()
	demuxer.OnFrame = func(cid mpeg2.TS_STREAM_TYPE, frame []byte, pts uint64, dts uint64) {
		switch cid {
		case mpeg2.TS_STREAM_H264:
			if !hasVideo {
				vtid = muxer.AddVideoTrack(mp4.MP4_CODEC_H264)
				hasVideo = true
			}
			err = muxer.Write(vtid, frame, pts, dts)
			if err != nil {
				return
				panic(err)
			}

		case mpeg2.TS_STREAM_AAC:
			if !hasAudio {
				atid = muxer.AddAudioTrack(mp4.MP4_CODEC_AAC)
				hasAudio = true
			}
			err = muxer.Write(atid, frame, pts, dts)
			if err != nil {
				return
				panic(err)
			}

		case mpeg2.TS_STREAM_AUDIO_MPEG1, mpeg2.TS_STREAM_AUDIO_MPEG2:
			if !hasAudio {
				atid = muxer.AddAudioTrack(mp4.MP4_CODEC_MP3)
				hasAudio = true
			}
			err = muxer.Write(atid, frame, pts, dts)
			if err != nil {
				return
				panic(err)
			}

		}

	}

	err = demuxer.Input(r)
	if err != nil {
		return err
	}

	err = muxer.WriteTrailer()
	if err != nil {
		return err
	}

	return nil
}
