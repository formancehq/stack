// Code generated by go-swagger; DO NOT EDIT.

package credit_transfers

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

// NewPostV1CreditTransfersParams creates a new PostV1CreditTransfersParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPostV1CreditTransfersParams() *PostV1CreditTransfersParams {
	return &PostV1CreditTransfersParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPostV1CreditTransfersParamsWithTimeout creates a new PostV1CreditTransfersParams object
// with the ability to set a timeout on a request.
func NewPostV1CreditTransfersParamsWithTimeout(timeout time.Duration) *PostV1CreditTransfersParams {
	return &PostV1CreditTransfersParams{
		timeout: timeout,
	}
}

// NewPostV1CreditTransfersParamsWithContext creates a new PostV1CreditTransfersParams object
// with the ability to set a context for a request.
func NewPostV1CreditTransfersParamsWithContext(ctx context.Context) *PostV1CreditTransfersParams {
	return &PostV1CreditTransfersParams{
		Context: ctx,
	}
}

// NewPostV1CreditTransfersParamsWithHTTPClient creates a new PostV1CreditTransfersParams object
// with the ability to set a custom HTTPClient for a request.
func NewPostV1CreditTransfersParamsWithHTTPClient(client *http.Client) *PostV1CreditTransfersParams {
	return &PostV1CreditTransfersParams{
		HTTPClient: client,
	}
}

/*
PostV1CreditTransfersParams contains all the parameters to send to the API endpoint

	for the post v1 credit transfers operation.

	Typically these are written to a http.Request.
*/
type PostV1CreditTransfersParams struct {

	/* CreditTransfer.

	   CreditTransfer
	*/
	CreditTransfer *models.CreatePaymentRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the post v1 credit transfers params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostV1CreditTransfersParams) WithDefaults() *PostV1CreditTransfersParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the post v1 credit transfers params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostV1CreditTransfersParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the post v1 credit transfers params
func (o *PostV1CreditTransfersParams) WithTimeout(timeout time.Duration) *PostV1CreditTransfersParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post v1 credit transfers params
func (o *PostV1CreditTransfersParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post v1 credit transfers params
func (o *PostV1CreditTransfersParams) WithContext(ctx context.Context) *PostV1CreditTransfersParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post v1 credit transfers params
func (o *PostV1CreditTransfersParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post v1 credit transfers params
func (o *PostV1CreditTransfersParams) WithHTTPClient(client *http.Client) *PostV1CreditTransfersParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post v1 credit transfers params
func (o *PostV1CreditTransfersParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCreditTransfer adds the creditTransfer to the post v1 credit transfers params
func (o *PostV1CreditTransfersParams) WithCreditTransfer(creditTransfer *models.CreatePaymentRequest) *PostV1CreditTransfersParams {
	o.SetCreditTransfer(creditTransfer)
	return o
}

// SetCreditTransfer adds the creditTransfer to the post v1 credit transfers params
func (o *PostV1CreditTransfersParams) SetCreditTransfer(creditTransfer *models.CreatePaymentRequest) {
	o.CreditTransfer = creditTransfer
}

// WriteToRequest writes these params to a swagger request
func (o *PostV1CreditTransfersParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.CreditTransfer != nil {
		if err := r.SetBodyParam(o.CreditTransfer); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
