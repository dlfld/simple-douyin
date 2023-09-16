package videoproducer

import (
	"context"
	"fmt"
	"io"

	"github.com/douyin/common/kafkaLog"
	"github.com/segmentio/kafka-go"
)

var kw *kafka.Writer

func init() {
	// 初始化kafka
	kw = &kafka.Writer{
		Addr:                   kafka.TCP(kafkaLog.KafkaAddr),
		Balancer:               &kafka.Hash{},
		Topic:                  "video",
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafka.RequireNone,
		Compression:            kafka.Snappy,
		BatchBytes:             1024 * 1024 * 50,
	}
}

func WriteVideoToKafka(reader io.Reader, dataLen, userID int64, title string) error {
	data := make([]byte, dataLen)
	reader.Read(data)
	fmt.Println("write video publish task to kafka: ", userID, title, dataLen)
	err := kw.WriteMessages(context.Background(),
		kafka.Message{
			Headers: []kafka.Header{{Key: "UserID", Value: []byte(fmt.Sprintf("%d", userID))}, {Key: "Title", Value: []byte(title)}},
			Value:   []byte(data),
		},
	)
	return err
}
