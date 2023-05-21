package jobsrepo

import (
	"fmt"

	"github.com/gerladeno/chat-service/internal/store"
)

//go:generate options-gen -out-filename=repo_options.gen.go -from-struct=Options
type Options struct {
	db *store.Database `option:"mandatory" validate:"required"`
}

type Repo struct {
	Options
}

func New(opts Options) (*Repo, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating jobs repo options: %v", err)
	}
	return &Repo{opts}, nil
}
