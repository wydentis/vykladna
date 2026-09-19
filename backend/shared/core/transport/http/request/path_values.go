package core_http_request

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

func GetStringPathValue(r *http.Request, key string) (string, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return "", fmt.Errorf("no key '%s' in path: %w", key, core_errors.ErrInvalidArgument)
	}

	return pathValue, nil
}

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	idString, err := GetStringPathValue(r, key)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("get from path error: %w", err)
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid uuid: %w", err)
	}

	return id, nil
}
