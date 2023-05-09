package msgproducer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/gerladeno/chat-service/internal/types"
)

type Message struct {
	ID         types.MessageID `json:"id"`
	ChatID     types.ChatID    `json:"chatId"` //nolint:tagliatelle
	Body       string          `json:"body"`
	FromClient bool            `json:"fromClient"`
}

func (s *Service) ProduceMessage(ctx context.Context, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshalling message: %v", err)
	}
	cipherText := data
	if s.cipher != nil {
		nonce, err := s.nonceFactory(s.cipher.NonceSize())
		if err != nil {
			return fmt.Errorf("deriving nonce: %v", err)
		}
		cipherText = s.cipher.Seal(nonce, nonce, data, nil)
	}
	if err = s.wr.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.ChatID.String()),
		Value: cipherText,
		Time:  time.Now(),
	}); err != nil {
		return fmt.Errorf("writing message to kafka: %v", err)
	}
	return nil
}

func (s *Service) Close() error {
	return s.wr.Close()
}
