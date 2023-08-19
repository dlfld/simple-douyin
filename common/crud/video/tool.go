package video

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
)

func GetSnapshot(videoPath, snapshotPath string, frameNum int) (ImagePath string, err error) {

	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()

	if err != nil {
		log.Fatal("生成略缩图失败")
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		//logger.Errorln("生成缩略图失败：", err)
		return "", err
	}

	err = imaging.Save(img, snapshotPath+".png")
	if err != nil {
		//logger.Errorln("生成缩略图失败：", err)
		return "", err
	}

	imgPath := snapshotPath + ".png"

	return imgPath, nil
}

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
