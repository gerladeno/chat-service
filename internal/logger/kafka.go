package logger

import (
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var _ kafka.Logger = (*KafkaAdapted)(nil)

type KafkaAdapted struct {
	log *zap.Logger
}

func (k KafkaAdapted) Printf(s string, i ...interface{}) {
	k.log.Sugar().Info(s, i)
}

func NewKafkaAdapted() *KafkaAdapted {
	return &KafkaAdapted{
		log: zap.L(),
	}
}

func (k KafkaAdapted) WithServiceName(name string) *KafkaAdapted {
	return &KafkaAdapted{
		log: k.log.With(zap.String("service_name", name)),
	}
}

func (k KafkaAdapted) ForErrors() *KafkaAdapted {
	return &KafkaAdapted{
		log: k.log.With(zap.String("no_clue", "error")),
	}
}
