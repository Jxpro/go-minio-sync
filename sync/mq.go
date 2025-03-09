package sync

import (
	rocketmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"go-minio-sync/config"
	"time"
)

type RocketMQ struct {
	Producer rocketmq.Producer
	Consumer rocketmq.SimpleConsumer
}

func NewRocketInstance(cfg *config.Config) (*RocketMQ, error) {
	producer, err := rocketmq.NewProducer(&rocketmq.Config{
		Endpoint: cfg.MQ.Endpoint,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    cfg.MQ.AccessKey,
			AccessSecret: cfg.MQ.SecretKey,
		},
	},
		// 实际使用时，message 中指定的 topic 会覆盖此处的配置
		rocketmq.WithTopics(cfg.MQ.Topic),
	)
	if err != nil {
		return nil, err
	}
	err = producer.Start()
	if err != nil {
		return nil, err
	}

	consumer, err := rocketmq.NewSimpleConsumer(&rocketmq.Config{
		Endpoint:      cfg.MQ.Endpoint,
		ConsumerGroup: cfg.MQ.ConsumerGroup,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    cfg.MQ.AccessKey,
			AccessSecret: cfg.MQ.SecretKey,
		},
	},
		rocketmq.WithAwaitDuration(time.Second*time.Duration(cfg.MQ.AwaitDuration)),
		rocketmq.WithSubscriptionExpressions(map[string]*rocketmq.FilterExpression{
			cfg.MQ.Topic: rocketmq.SUB_ALL,
		}),
	)
	if err != nil {
		return nil, err
	}
	err = consumer.Start()
	if err != nil {
		return nil, err
	}

	return &RocketMQ{
		Producer: producer,
		Consumer: consumer,
	}, nil
}

func (m *RocketMQ) Shutdown() error {
	if err := m.Producer.GracefulStop(); err != nil {
		return err
	}
	if err := m.Consumer.GracefulStop(); err != nil {
		return err
	}
	return nil
}
