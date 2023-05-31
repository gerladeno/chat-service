package canreceiveproblems

import (
	"context"
	"fmt"

	"github.com/gerladeno/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=canreceiveproblemsmock

type managerLoadService interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type managerPool interface {
	Contains(ctx context.Context, managerID types.UserID) (bool, error)
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

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("request validation: %v", err)
	}

	contains, err := u.managerPool.Contains(ctx, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("check if manager pool contains specified manager id already: %v", err)
	}
	if contains {
		return Response{Result: false}, nil
	}

	can, err := u.managerLoadService.CanManagerTakeProblem(ctx, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("check if managers capacity is already exceeded: %v", err)
	}
	return Response{Result: can}, nil
}
