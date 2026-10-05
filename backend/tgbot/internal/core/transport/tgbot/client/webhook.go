package core_tgbot_client

import "context"

type WebhookParams struct {
	URL                string   `json:"url"`
	SecretToken        string   `json:"secret_token,omitempty"`
	AllowedUpdates     []string `json:"allowed_updates,omitempty"`
	DropPendingUpdates bool     `json:"drop_pending_updates,omitempty"`
}

type WebhookInfo struct {
	URL                string   `json:"url"`
	PendingUpdateCount int      `json:"pending_update_count"`
	LastErrorDate      int64    `json:"last_error_date"`
	LastErrorMessage   string   `json:"last_error_message"`
	AllowedUpdates     []string `json:"allowed_updates"`
}

func (c *Client) SetWebhook(ctx context.Context, p WebhookParams) error {
	return c.call(ctx, "setWebhook", p, nil)
}

func (c *Client) DeleteWebhook(ctx context.Context, dropPending bool) error {
	return c.call(ctx, "deleteWebhook", map[string]bool{"drop_pending_updates": dropPending}, nil)
}

func (c *Client) WebhookInfo(ctx context.Context) (WebhookInfo, error) {
	var info WebhookInfo
	err := c.call(ctx, "getWebhookInfo", struct{}{}, &info)
	return info, err
}
