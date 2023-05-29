package getchathistory

import (
	"context"
	"errors"
	"fmt"

	"github.com/gerladeno/chat-service/internal/cursor"
	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	"github.com/gerladeno/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=getchathistorymocks

type messagesRepository interface {
	GetManagerChatMessages(
		ctx context.Context,
		chatID types.ChatID,
		managerID types.UserID,
		pageSize int,
		cursor *messagesrepo.Cursor,
	) ([]messagesrepo.Message, *messagesrepo.Cursor, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("validating get chat history usecase options: %v", err)
	}
	return UseCase{opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, ErrInvalidRequest
	}
	var cursorParam *messagesrepo.Cursor
	if req.Cursor != "" {
		var reqCursor messagesrepo.Cursor
		if err := cursor.Decode(req.Cursor, &reqCursor); err != nil {
			return Response{}, ErrInvalidCursor
		}
		cursorParam = &reqCursor
	}

	messages, respCursor, err := u.msgRepo.GetManagerChatMessages(ctx, req.ChatID, req.ManagerID, req.PageSize, cursorParam)
	switch {
	case errors.Is(err, messagesrepo.ErrInvalidCursor):
		return Response{}, ErrInvalidCursor
	case err != nil:
		return Response{}, fmt.Errorf("GetClientChatMessages: %v", err)
	}

	resp := Response{}
	if respCursor != nil {
		resp.NextCursor, err = cursor.Encode(respCursor)
		if err != nil {
			return Response{}, fmt.Errorf("encoding next cursor: %v", err)
		}
	}
	for i := range messages {
		resp.Messages = append(resp.Messages, Message{
			ID:        messages[i].ID,
			AuthorID:  messages[i].AuthorID,
			Body:      messages[i].Body,
			CreatedAt: messages[i].CreatedAt,
		})
	}
	return resp, nil
}
