package video

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/snowflake"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"os"
)

//	GenSnowId
//
// @Description:	雪花算法生成分布式id
// @return string
// @return error
func GenSnowId() (string, error) {
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Generate a snowflake ID.
	id := node.Generate()
	//fmt.Printf("String ID: %s\n", id)
	return fmt.Sprintf("%s", id), nil
}

//	GetSnapshotImageBuffer
//
// @Description:	截取视频的第一帧作为视频的封面图片
// @param videoPath
// @param frameNum
// @return *bytes.Buffer
// @return error
func GetSnapshotImageBuffer(videoPath string, frameNum int) (*bytes.Buffer, error) {
	//logger := zap.InitLogger()

	buf := bytes.NewBuffer(nil)

	err := ffmpeg.Input(videoPath).Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()

	if err != nil {
		//logger.Errorln("生成缩略图失败：", err)
		return nil, err
	}

	return buf, nil
}
