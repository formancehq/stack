// Code generated by go-swagger; DO NOT EDIT.

package external_accounts

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

// NewPutV1ExternalAccountsIDParams creates a new PutV1ExternalAccountsIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPutV1ExternalAccountsIDParams() *PutV1ExternalAccountsIDParams {
	return &PutV1ExternalAccountsIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPutV1ExternalAccountsIDParamsWithTimeout creates a new PutV1ExternalAccountsIDParams object
// with the ability to set a timeout on a request.
func NewPutV1ExternalAccountsIDParamsWithTimeout(timeout time.Duration) *PutV1ExternalAccountsIDParams {
	return &PutV1ExternalAccountsIDParams{
		timeout: timeout,
	}
}

// NewPutV1ExternalAccountsIDParamsWithContext creates a new PutV1ExternalAccountsIDParams object
// with the ability to set a context for a request.
func NewPutV1ExternalAccountsIDParamsWithContext(ctx context.Context) *PutV1ExternalAccountsIDParams {
	return &PutV1ExternalAccountsIDParams{
		Context: ctx,
	}
}

// NewPutV1ExternalAccountsIDParamsWithHTTPClient creates a new PutV1ExternalAccountsIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewPutV1ExternalAccountsIDParamsWithHTTPClient(client *http.Client) *PutV1ExternalAccountsIDParams {
	return &PutV1ExternalAccountsIDParams{
		HTTPClient: client,
	}
}

/*
PutV1ExternalAccountsIDParams contains all the parameters to send to the API endpoint

	for the put v1 external accounts ID operation.

	Typically these are written to a http.Request.
*/
type PutV1ExternalAccountsIDParams struct {

	/* ExternalAccount.

	   ExternalAccount
	*/
	ExternalAccount *models.UpdateExternalAccountRequest

	/* ID.

	   External account ID
	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the put v1 external accounts ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutV1ExternalAccountsIDParams) WithDefaults() *PutV1ExternalAccountsIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the put v1 external accounts ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutV1ExternalAccountsIDParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the put v1 external accounts ID params
func (o *PutV1ExternalAccountsIDParams) WithTimeout(timeout time.Duration) *PutV1ExternalAccountsIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the put v1 external accounts ID params
func (o *PutV1ExternalAccountsIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the put v1 external accounts ID params
func (o *PutV1ExternalAccountsIDParams) WithContext(ctx context.Context) *PutV1ExternalAccountsIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the put v1 external accounts ID params
func (o *PutV1ExternalAccountsIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the put v1 external accounts ID params
func (o *PutV1ExternalAccountsIDParams) WithHTTPClient(client *http.Client) *PutV1ExternalAccountsIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the put v1 external accounts ID params
func (o *PutV1ExternalAccountsIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithExternalAccount adds the externalAccount to the put v1 external accounts ID params
func (o *PutV1ExternalAccountsIDParams) WithExternalAccount(externalAccount *models.UpdateExternalAccountRequest) *PutV1ExternalAccountsIDParams {
	o.SetExternalAccount(externalAccount)
	return o
}

// SetExternalAccount adds the externalAccount to the put v1 external accounts ID params
func (o *PutV1ExternalAccountsIDParams) SetExternalAccount(externalAccount *models.UpdateExternalAccountRequest) {
	o.ExternalAccount = externalAccount
}

// WithID adds the id to the put v1 external accounts ID params
func (o *PutV1ExternalAccountsIDParams) WithID(id string) *PutV1ExternalAccountsIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the put v1 external accounts ID params
func (o *PutV1ExternalAccountsIDParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *PutV1ExternalAccountsIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.ExternalAccount != nil {
		if err := r.SetBodyParam(o.ExternalAccount); err != nil {
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