package models

type PSPWebhookConfig struct {
	Name    string `json:"name"`
	URLPath string `json:"urlPath"`
}

type WebhookConfig struct {
	Name        string      `json:"name"`
	ConnectorID ConnectorID `json:"connectorID"`
	URLPath     string      `json:"urlPath"`
}

type PSPWebhook struct {
	QueryValues map[string][]string `json:"queryValues"`
	Headers     map[string][]string `json:"headers"`
	Body        []byte              `json:"payload"`
}

type Webhook struct {
	ID          string              `json:"id"`
	ConnectorID ConnectorID         `json:"connectorID"`
	QueryValues map[string][]string `json:"queryValues"`
	Headers     map[string][]string `json:"headers"`
	Body        []byte              `json:"payload"`
}
