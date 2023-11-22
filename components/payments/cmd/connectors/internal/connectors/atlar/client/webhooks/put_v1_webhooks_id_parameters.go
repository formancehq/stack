// Code generated by go-swagger; DO NOT EDIT.

package webhooks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/models"
)

// NewPutV1WebhooksIDParams creates a new PutV1WebhooksIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPutV1WebhooksIDParams() *PutV1WebhooksIDParams {
	return &PutV1WebhooksIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPutV1WebhooksIDParamsWithTimeout creates a new PutV1WebhooksIDParams object
// with the ability to set a timeout on a request.
func NewPutV1WebhooksIDParamsWithTimeout(timeout time.Duration) *PutV1WebhooksIDParams {
	return &PutV1WebhooksIDParams{
		timeout: timeout,
	}
}

// NewPutV1WebhooksIDParamsWithContext creates a new PutV1WebhooksIDParams object
// with the ability to set a context for a request.
func NewPutV1WebhooksIDParamsWithContext(ctx context.Context) *PutV1WebhooksIDParams {
	return &PutV1WebhooksIDParams{
		Context: ctx,
	}
}

// NewPutV1WebhooksIDParamsWithHTTPClient creates a new PutV1WebhooksIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewPutV1WebhooksIDParamsWithHTTPClient(client *http.Client) *PutV1WebhooksIDParams {
	return &PutV1WebhooksIDParams{
		HTTPClient: client,
	}
}

/*
PutV1WebhooksIDParams contains all the parameters to send to the API endpoint

	for the put v1 webhooks ID operation.

	Typically these are written to a http.Request.
*/
type PutV1WebhooksIDParams struct {

	/* Webhook.

	   Webhook
	*/
	Webhook *models.UpdateWebhookRequest

	/* ID.

	   Webhook ID
	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the put v1 webhooks ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutV1WebhooksIDParams) WithDefaults() *PutV1WebhooksIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the put v1 webhooks ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutV1WebhooksIDParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the put v1 webhooks ID params
func (o *PutV1WebhooksIDParams) WithTimeout(timeout time.Duration) *PutV1WebhooksIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the put v1 webhooks ID params
func (o *PutV1WebhooksIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the put v1 webhooks ID params
func (o *PutV1WebhooksIDParams) WithContext(ctx context.Context) *PutV1WebhooksIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the put v1 webhooks ID params
func (o *PutV1WebhooksIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the put v1 webhooks ID params
func (o *PutV1WebhooksIDParams) WithHTTPClient(client *http.Client) *PutV1WebhooksIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the put v1 webhooks ID params
func (o *PutV1WebhooksIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithWebhook adds the webhook to the put v1 webhooks ID params
func (o *PutV1WebhooksIDParams) WithWebhook(webhook *models.UpdateWebhookRequest) *PutV1WebhooksIDParams {
	o.SetWebhook(webhook)
	return o
}

// SetWebhook adds the webhook to the put v1 webhooks ID params
func (o *PutV1WebhooksIDParams) SetWebhook(webhook *models.UpdateWebhookRequest) {
	o.Webhook = webhook
}

// WithID adds the id to the put v1 webhooks ID params
func (o *PutV1WebhooksIDParams) WithID(id string) *PutV1WebhooksIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the put v1 webhooks ID params
func (o *PutV1WebhooksIDParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *PutV1WebhooksIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Webhook != nil {
		if err := r.SetBodyParam(o.Webhook); err != nil {
			return err
		}
	}

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
