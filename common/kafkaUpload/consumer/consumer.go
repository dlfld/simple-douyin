package videoconsumer

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/douyin/common/kafkaLog"
	"github.com/douyin/common/kafkaLog/productor"
	"github.com/douyin/common/utils"
	"github.com/segmentio/kafka-go"
)

var kr *kafka.Reader
var logger *productor.LogCollector

func init() {
	kr = kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:        []string{fmt.Sprintf(kafkaLog.KafkaAddr)},
			GroupID:        "videoPublisher",
			Topic:          "video",
			CommitInterval: time.Second * 2,
			StartOffset:    kafka.FirstOffset,
			MinBytes:       1024,
			MaxBytes:       100 * 1024 * 1024,
		},
	)
	var err error
	logger, err = productor.NewLogCollector("videoPublisher")
	if err != nil {
		panic(err)
	}
}

// HandleVideoFromKafka 从kafka中读取待处理的视频
func HandleVideoFromKafka() error {
	for {
		msg, err := kr.ReadMessage(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("failed read kafka, err=%s", err.Error()))
			return err
		}
		userId, err := strconv.ParseInt(string(msg.Headers[0].Value), 10, 64)
		if err != nil {
			logger.Error(fmt.Sprintf("failed parse userId, err=%s", err.Error()))

			return err
		}
		title := string(msg.Headers[1].Value)
		reader := bytes.NewReader(msg.Value)
		dataLen := int64(len(msg.Value))
		// save tmp video
		videoPath := fmt.Sprintf("/tmp/%d-%d.mp4", userId, time.Now().Unix())
		sv, err := os.Create(videoPath)
		if err != nil {
			logger.Error(fmt.Sprintf("failed create file, err=%s", err.Error()))
			return err
		}
		_, err = sv.Write(msg.Value)
		if err != nil {
			logger.Error(fmt.Sprintf("failed write file, err=%s", err.Error()))
			return err
		}
		// add watermark
		watermarkPath, err := utils.GetWaterMark(int(userId), reader)
		if err != nil {
			logger.Error(fmt.Sprintf("failed get watermark, err=%s", err.Error()))
			return err

		}
		reader, err = utils.AddWatermarkToVideo(watermarkPath, videoPath)
		if err != nil {
			logger.Error(fmt.Sprintf("failed add watermark, err=%s", err.Error()))
			return err
		}
		// upload video
		for i := 0; i < 3; i++ {
			err = utils.UploadVideo(reader, dataLen, userId, title)
			if err != nil {
				logger.Error(fmt.Sprintf("retry [%d] failed upload video, err=%s", i, err.Error()))
				return err
			}

		}
		// writeoffset
		kr.CommitMessages(context.Background(), msg)

	}
}
