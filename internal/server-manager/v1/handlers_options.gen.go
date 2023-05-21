// Code generated by options-gen. DO NOT EDIT.
package managerv1

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"go.uber.org/zap"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	logger *zap.Logger,
	canReceiveProblemsUseCase canReceiveProblemsUseCase,
	freeHandsUseCase freeHandsUseCase,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.logger = logger
	o.canReceiveProblemsUseCase = canReceiveProblemsUseCase
	o.freeHandsUseCase = freeHandsUseCase

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("logger", _validate_Options_logger(o)))
	errs.Add(errors461e464ebed9.NewValidationError("canReceiveProblemsUseCase", _validate_Options_canReceiveProblemsUseCase(o)))
	errs.Add(errors461e464ebed9.NewValidationError("freeHandsUseCase", _validate_Options_freeHandsUseCase(o)))
	return errs.AsError()
}

func _validate_Options_logger(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.logger, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `logger` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_canReceiveProblemsUseCase(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.canReceiveProblemsUseCase, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `canReceiveProblemsUseCase` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_freeHandsUseCase(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.freeHandsUseCase, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `freeHandsUseCase` did not pass the test: %w", err)
	}
	return nil
}
