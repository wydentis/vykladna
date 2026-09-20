package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

type hasCustomValidation interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	v, ok := dest.(hasCustomValidation)
	if ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validate request: %v: %w", err, core_errors.ErrInvalidArgument)
		}
	}

	return nil
}
