// Code generated by options-gen. DO NOT EDIT.
package managerworkspace

import (
	fmt461e464ebed9 "fmt"

	"github.com/gerladeno/chat-service/internal/types"
	apimanagerv1 "github.com/gerladeno/chat-service/tests/e2e/api/manager/v1"
	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	id types.UserID,
	token string,
	api *apimanagerv1.ClientWithResponses,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.id = id
	o.token = token
	o.api = api

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("id", _validate_Options_id(o)))
	errs.Add(errors461e464ebed9.NewValidationError("token", _validate_Options_token(o)))
	errs.Add(errors461e464ebed9.NewValidationError("api", _validate_Options_api(o)))
	return errs.AsError()
}

func _validate_Options_id(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.id, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `id` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_token(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.token, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `token` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_api(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.api, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `api` did not pass the test: %w", err)
	}
	return nil
}
