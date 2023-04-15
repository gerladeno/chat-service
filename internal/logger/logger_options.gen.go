// Code generated by options-gen. DO NOT EDIT.
package logger

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	level string,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.level = level

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func WithProductionMode(opt bool) OptOptionsSetter {
	return func(o *Options) {
		o.productionMode = opt
	}
}

func WithSentryDSN(opt string) OptOptionsSetter {
	return func(o *Options) {
		o.sentryDSN = opt
	}
}

func WithEnv(opt string) OptOptionsSetter {
	return func(o *Options) {
		o.env = opt
	}
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("level", _validate_Options_level(o)))
	return errs.AsError()
}

func _validate_Options_level(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.level, "required,oneof=debug info warn error"); err != nil {
		return fmt461e464ebed9.Errorf("field `level` did not pass the test: %w", err)
	}
	return nil
}
