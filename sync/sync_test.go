package sync

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	rocketmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "go-minio-sync/config"
)

var cfg = Config{
	MQ: MQConfig{
		Topic:         "ordTest2",
		Endpoint:      "172.20.165.191:8081",
		ConsumerGroup: "test-group",
		AwaitDuration: 5,
	},
}

func TestSaveAndLoadState(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "temp1g_*")
	require.NoError(t, err, "文件创建失败")
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	const oneGB = int64(1024 * 1024 * 1024)
	err = tmpFile.Truncate(oneGB)
	assert.NoError(t, err, "文件截断失败")

	state := &State{
		FilePath:    tmpFile.Name(),
		UploadID:    "12345",
		FileSize:    123456789,
		TrunkSize:   5 * 1024 * 1024, // 5MB
		TrunkLength: 123456789 / (5 * 1024 * 1024),
	}

	err = state.Save()
	require.NoError(t, err, "文件保存失败")

	data, err := os.ReadFile(tmpFile.Name() + ".upload.state")
	require.NoError(t, err, "文件读取失败")
	loaded, err := LoadState(data)
	require.NoError(t, err, "文件解析失败")

	assert.Equal(t, state.UploadID, loaded.UploadID, "UploadID 不一致")
	assert.Equal(t, state.TrunkLength, loaded.TrunkLength, "分块大小不一致")
	assert.Equal(t, state.FileSize, loaded.FileSize, "文件大小不一致")
}

func TestMessageQueueInit(t *testing.T) {
	mq, err := NewRocketInstance(&cfg)
	require.NoError(t, err, "初始化 RocketMQ 客户端失败")

	err = mq.Shutdown()
	require.NoError(t, err, "关闭 RocketMQ 客户端失败")
}

func TestMessageQueueSend(t *testing.T) {
	mq, err := NewRocketInstance(&cfg)
	require.NoError(t, err, "初始化 RocketMQ 客户端失败")

	t.Log("开始发送消息")
	for i := 0; i < 100; i++ {
		msg := &rocketmq.Message{
			Topic: cfg.MQ.Topic,
			Body:  []byte(fmt.Sprintf("this is a message %d:  in group-%d", i/10, i%10)),
		}

		// 发送同步消息
		msg.SetMessageGroup(fmt.Sprintf("group-%d", i%10))

		resp, err := mq.Producer.Send(context.TODO(), msg)
		assert.NoError(t, err, "发送消息失败")
		t.Logf("发送消息成功: %#v", resp[0])
	}

	err = mq.Shutdown()
	require.NoError(t, err, "关闭 RocketMQ 客户端失败")
}

func TestMessageQueueReceive(t *testing.T) {
	mq, err := NewRocketInstance(&cfg)
	require.NoError(t, err, "初始化 RocketMQ 客户端失败")

	t.Log("开始接收消息")
	for {
		mvs, err := mq.Consumer.Receive(context.TODO(), 1000, time.Second*10)
		if err != nil && strings.Contains(err.Error(), "no new message") {
			break
		}

		// 响应消息
		for _, mv := range mvs {
			_ = mq.Consumer.Ack(context.TODO(), mv)
			t.Log(string(mv.GetBody()))
		}
	}

	err = mq.Shutdown()
	require.NoError(t, err, "关闭 RocketMQ 客户端失败")
}
