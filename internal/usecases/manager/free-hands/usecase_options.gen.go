// Code generated by options-gen. DO NOT EDIT.
package freehands

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	managerLoadService managerLoadService,
	managerPool managerPool,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.managerLoadService = managerLoadService
	o.managerPool = managerPool

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("managerLoadService", _validate_Options_managerLoadService(o)))
	errs.Add(errors461e464ebed9.NewValidationError("managerPool", _validate_Options_managerPool(o)))
	return errs.AsError()
}

func _validate_Options_managerLoadService(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.managerLoadService, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `managerLoadService` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_managerPool(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.managerPool, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `managerPool` did not pass the test: %w", err)
	}
	return nil
}
