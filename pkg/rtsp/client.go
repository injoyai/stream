package rtsp

import (
	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/description"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/injoyai/ios"
	"github.com/pion/rtp"
	"io"
	"sync"
)

type (
	Media  = description.Media
	Format = format.Format
	Packet = rtp.Packet
)

func Dial(url string) (*Client, error) {
	c, err := New(url)
	if err != nil {
		return nil, err
	}
	return c, c.Connect()
}

func New(url string) (*Client, error) {
	c := &gortsplib.Client{}
	u, err := base.ParseURL(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: c,
		URL:    u,
	}, nil
}

type Client struct {
	*gortsplib.Client
	URL  *base.URL
	once sync.Once
}

func (this *Client) Connect() error {
	// 连接到 RTSP 服务端
	err := this.Client.Start(this.URL.Scheme, this.URL.Host)
	if err != nil {
		return err
	}
	// 获取描述选项啥的
	desc, _, err := this.Client.Describe(this.URL)
	if err != nil {
		return err
	}
	// 拉取啥类型的流
	err = this.Client.SetupAll(desc.BaseURL, desc.Medias)
	if err != nil {
		return err
	}
	// 开始拉取
	_, err = this.Client.Play(nil)
	return err
}

// OnPacket 处理读取到的包,需要在Connect之后
func (this *Client) OnPacket(f func(media *Media, format Format, packet *Packet)) {
	this.Client.OnPacketRTPAny(f)
}

func (this *Client) Close() error {
	this.Client.Close()
	return nil
}

func (this *Client) Done() <-chan struct{} {
	this.Client.Wait()
	done := make(chan struct{})
	close(done)
	return done
}

// MReadCloser 转成ios.MReadCloser,方便复制数据
func (this *Client) MReadCloser() ios.MReadCloser {
	c := make(chan []byte)
	this.OnPacket(func(media *Media, format Format, packet *Packet) {
		c <- packet.Payload
	})
	return struct {
		ios.MReader
		io.Closer
	}{
		MReader: ios.MReadFunc(func() ([]byte, error) {
			bs, ok := <-c
			if !ok {
				return nil, io.EOF
			}
			return bs, nil
		}),
		Closer: this,
	}
}
