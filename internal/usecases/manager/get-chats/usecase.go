package getchats

import (
	"context"
	"fmt"

	chatsrepo "github.com/gerladeno/chat-service/internal/repositories/chats"
	"github.com/gerladeno/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=getchatsmock

type chatsRepo interface {
	GetChatsForManager(ctx context.Context, managerID types.UserID) ([]chatsrepo.Chat, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	chatsRepo chatsRepo `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("validate get chats usecase options: %v", err)
	}
	return UseCase{Options: opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("invalid getchats request: %v", err)
	}
	chats, err := u.chatsRepo.GetChatsForManager(ctx, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("GetChatsForManager: %v", err)
	}
	result := make([]Chat, 0, len(chats))
	for _, chat := range chats {
		result = append(result, Chat{
			ID:       chat.ID,
			ClientID: chat.ClientID,
		})
	}
	return Response{Chats: result}, nil
}
