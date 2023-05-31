package logger

import (
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var _ kafka.Logger = (*KafkaAdapted)(nil)

type KafkaAdapted struct {
	forErrors bool
	log       *zap.Logger
}

func (k *KafkaAdapted) Printf(s string, args ...interface{}) {
	if k.forErrors {
		k.log.Sugar().Errorf(s, args...)
	} else {
		k.log.Sugar().Debugf(s, args...)
	}
}

func NewKafkaAdapted() *KafkaAdapted {
	return &KafkaAdapted{
		log: zap.L(),
	}
}

func (k *KafkaAdapted) WithServiceName(name string) *KafkaAdapted {
	k.log = k.log.Named(name)
	return k
}

func (k *KafkaAdapted) ForErrors() *KafkaAdapted {
	k.forErrors = true
	return k
}
