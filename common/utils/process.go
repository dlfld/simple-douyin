package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/douyin/common/crud"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

func GetWaterMark(userID int, reader io.Reader) (filepath string, err error) {
	fontBytes, err := os.ReadFile("/app/static/simkai.ttf")
	filepath = fmt.Sprintf("/tmp/%d-%d.png", userID, time.Now().Unix())
	if err != nil {
		return
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return
	}
	fontSize := 60
	imgWidth := 800
	imgHeight := 80
	user, err := crud.GetUserInfo(strconv.Itoa(userID))
	if err != nil {
		return
	}
	text := user.UserName
	textColor := color.RGBA{R: 200, G: 200, B: 200, A: 200}
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.Transparent}, image.Point{}, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(float64(fontSize))
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(textColor))
	textX := 10
	textY := fontSize
	pt := freetype.Pt(textX, textY)
	_, err = c.DrawString(text, pt)
	if err != nil {
		return
	}
	var file *os.File
	file, err = os.Create(filepath)
	if err != nil {
		return
	}
	err = png.Encode(file, img)
	if err != nil {
		return
	}
	return
}

func AddWatermarkToVideo(waterMarkPath string, videoPath string) (reader *bytes.Reader, err error) {
	cmdArgs := []string{
		"-i", videoPath,
		"-i", waterMarkPath,
		"-filter_complex", "[0:v][1:v]overlay=10:10",
		"-f", "matroska", "-",
	}
	cmd := exec.Command("ffmpeg", cmdArgs...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	// remove temp file
	os.Remove(waterMarkPath)
	os.Remove(videoPath)
	reader = bytes.NewReader(buf.Bytes())
	return
}
