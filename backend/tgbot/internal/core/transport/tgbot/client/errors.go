package core_tgbot_client

import "fmt"

type APIError struct {
	Method      string
	Status      int
	Description string
	RetryAfter  int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("telegram %s: status %d: %s", e.Method, e.Status, e.Description)
}
