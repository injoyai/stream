package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/injoyai/goutil/frame/in/v3"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/injoyai/goutil/frame/mux"
	"github.com/injoyai/logs"
)

func pushRTSPStream(rtspURL string) error {
	cmd := exec.Command("ffmpeg",
		"-rtsp_transport", "tcp",
		"-i", rtspURL,
		"-c:v", "mpeg1video",
		"-b:v", "800k", // 增加码率，提升画面质量
		"-r", "30", // 帧率 30fps
		//"-g", "15", // 每秒一个关键帧（关键帧间隔=帧率）
		"-bf", "0", // 禁用 B 帧，减少复杂性

		"-c:a", "mp2",
		"-ar", "44100",
		"-ac", "1",
		"-s", "680x480",
		"-f", "mpegts",
		"http://127.0.0.1:8080/push",
	)
	return cmd.Run()
}

func pushFile() error {
	f, err := os.Open("F:\\test\\out.mts")
	if err != nil {
		return err
	}
	defer f.Close()
	for {
		buf := make([]byte, 1024*8)
		n, err := f.Read(buf)
		if err != nil {
			return err
		}
		sub.Publish(buf[:n])
		<-time.After(time.Millisecond * 10)
	}

}

// 拉取 RTSP 流数据
func pullRTSPStream(rtspURL string, sub *Subscribe) error {
	//return pushFile()
	//return pushRTSPStream(rtspURL)
	logs.Info("开始拉取 RTSP 流:", rtspURL)

	//cmd := exec.Command("ffmpeg",
	//	"-rtsp_transport", "tcp",
	//	"-i", rtspURL,
	//	"-f", "mpegts",
	//	"-codec:v", "mpeg1video",
	//	"-b:v", "800k",
	//	"-r", "25",
	//	"-bf", "0",
	//	"-codec:a", "mp2",
	//	"-ar", "44100",
	//	"-ac", "1",
	//	"-f", "mpeg1video",
	//	"-")

	//cmd := exec.Command("ffmpeg",
	//	"-rtsp_transport", "tcp",
	//	"-i", rtspURL,
	//	//"-f", "mpegts",
	//	"-c:v", "mpeg1video", //"libx264", // "mpeg1video", //"mjpeg",
	//	//"-b:v", "800k",
	//	"-r", "30",
	//	"-bf", "0",
	//	"-c:a", "mp2",
	//	"-ar", "44100",
	//	"-ac", "1",
	//	"-s", "680x480",
	//	"-f", "mpegts",
	//	"-")

	cmd := exec.Command("ffmpeg",
		"-hwaccel", "cuda", // 启用 CUDA 加速
		"-hide_banner",
		//"-f", "mpegts",
		"-i", "F:\\test\\x36xhzz.mp4", //rtspURL,
		"-c:v", "mpeg1video",
		"-b:v", "1500k", // 增加码率，提升画面质量
		"-r", "30", // 帧率 30fps
		"-g", "30", // 每秒一个关键帧（关键帧间隔=帧率）
		"-bf", "0", // 禁用 B 帧，减少复杂性

		"-c:a", "mp2",
		"-ar", "44100",
		"-ac", "1",
		"-s", "680x480",
		"-f", "mpegts",
		"-",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("获取标准输出失败: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("获取标准错误失败: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动 FFmpeg 失败: %v", err)
	}

	// 读取错误输出
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			logs.Debug("FFmpeg:", scanner.Text())
		}
	}()

	// 读取视频流
	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := stdout.Read(buffer)
			if err != nil {
				logs.Err("读取 FFmpeg 输出错误:", err)
				break
			}
			if n > 0 {
				//sub.Publish(append(buffer[:n], '\n'))
				sub.Publish(buffer[:n])
				//logs.Debug("推送视频数据:", n, "字节")
			}
		}
	}()

	logs.Info("FFmpeg 进程已启动")
	return cmd.Wait()
}

var sub = NewSubscribe()

func main() {
	logs.Info("启动服务器...")

	s := mux.New(mux.WithPort(8080))

	logs.Info("注册 WebSocket 处理器...")
	s.ALL("/ws", handlerWs)
	s.ALL("/push", handlerPush)

	logs.Info("设置静态文件服务...")
	s.Static("/", "./html")

	// 启动 RTSP 拉流
	go func() {
		rtspURL := "rtsp://admin:ASRDEO@192.168.10.47:554"
		logs.Info("开始 RTSP 拉流循环...")
		for {
			err := pullRTSPStream(rtspURL, sub)
			if err != nil {
				logs.Error("RTSP 拉流错误:", err)
				time.Sleep(time.Second * 3)
			}
		}
	}()

	logs.Info("服务器启动在 :8080")
	s.Run()
}

// 启动 WebSocket 服务器
func handlerWs(r *mux.Request) {
	logs.Info("新的 WebSocket 连接请求")

	up := websocket.Upgrader{
		ReadBufferSize:  1024 * 2,
		WriteBufferSize: 1024 * 64,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	ws, err := up.Upgrade(r.Writer, r.Request, nil)
	in.CheckErr(err)

	//ws := r.Websocket()
	defer ws.Close()

	logs.Info("WebSocket 连接已建立")

	go func() {
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				return
			}
		}
	}()

	c := sub.Subscribe(10)
	defer c.Close()

	logs.Info("开始推送视频流")

	for {
		select {
		case data, ok := <-c.C:
			if !ok {
				logs.Error("订阅通道已关闭")
				return
			}

			if err := ws.WriteMessage(websocket.BinaryMessage, data); err != nil {
				logs.Error("写入 WebSocket 失败:", err)
				return
			}

		case <-r.Context().Done():
			logs.Info("客户端断开连接")
			return
		}
	}
}

func handlerPush(r *mux.Request) {
	defer r.Body.Close()
	reader := bufio.NewReader(r.Body)
	for {
		bs, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		sub.Publish(bs)
	}
}
