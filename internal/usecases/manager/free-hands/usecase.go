package freehands

import (
	"context"
	"errors"
	"fmt"

	"github.com/gerladeno/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=freehandsmocks

var ErrManagerOverloaded = errors.New("manager overloaded")

type managerLoadService interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type managerPool interface {
	Put(ctx context.Context, managerID types.UserID) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	managerLoadService managerLoadService `option:"mandatory" validate:"required"`
	managerPool        managerPool        `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (*UseCase, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating can receive problems usecase options: %v", err)
	}
	return &UseCase{Options: opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("request validation: %v", err)
	}
	can, err := u.managerLoadService.CanManagerTakeProblem(ctx, req.ManagerID)
	if err != nil {
		return fmt.Errorf("check if managers capacity is already exceeded: %v", err)
	}
	if !can {
		return ErrManagerOverloaded
	}
	if err = u.managerPool.Put(ctx, req.ManagerID); err != nil {
		return fmt.Errorf("put manager into manager pool: %v", err)
	}
	return nil
}
