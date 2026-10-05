package core_tgbot_client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	http *http.Client
	base string
}

func New(cfg Config) *Client {
	return &Client{
		http: &http.Client{Timeout: cfg.Timeout},
		base: cfg.APIURL + "/bot" + cfg.Token,
	}
}

func (c *Client) call(ctx context.Context, method string, in, out any) error {
	body, err := json.Marshal(in)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/"+method, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		var ue *url.Error
		if errors.As(err, &ue) {
			err = ue.Err
		}
		return fmt.Errorf("telegram %s: %w", method, err)
	}
	defer res.Body.Close()

	var env struct {
		OK          bool            `json:"ok"`
		Description string          `json:"description"`
		Result      json.RawMessage `json:"result"`
		Parameters  struct {
			RetryAfter int `json:"retry_after"`
		} `json:"parameters"`
	}
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		return fmt.Errorf("telegram %s: status %d: undecodable response: %w", method, res.StatusCode, err)
	}
	if !env.OK {
		return &APIError{Method: method, Status: res.StatusCode, Description: env.Description, RetryAfter: env.Parameters.RetryAfter}
	}
	if out != nil {
		return json.Unmarshal(env.Result, out)
	}

	return nil
}
