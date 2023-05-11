package managerload

import (
	"context"
	"fmt"

	"github.com/gerladeno/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/manager_load_mock.gen.go -package=managerloadmock

type problemsRepository interface {
	GetManagerOpenProblemsCount(ctx context.Context, managerID types.UserID) (int, error)
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	maxProblemsAtTime int                `option:"mandatory" validate:"required,min=1,max=30"`
	problemsRepo      problemsRepository `option:"mandatory" validate:"required"`
}

type Service struct {
	Options
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate manager load options: %v", err)
	}
	return &Service{Options: opts}, nil
}
