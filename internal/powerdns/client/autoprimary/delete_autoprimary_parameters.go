// Code generated by go-swagger; DO NOT EDIT.

package autoprimary

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
)

// NewDeleteAutoprimaryParams creates a new DeleteAutoprimaryParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteAutoprimaryParams() *DeleteAutoprimaryParams {
	return &DeleteAutoprimaryParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteAutoprimaryParamsWithTimeout creates a new DeleteAutoprimaryParams object
// with the ability to set a timeout on a request.
func NewDeleteAutoprimaryParamsWithTimeout(timeout time.Duration) *DeleteAutoprimaryParams {
	return &DeleteAutoprimaryParams{
		timeout: timeout,
	}
}

// NewDeleteAutoprimaryParamsWithContext creates a new DeleteAutoprimaryParams object
// with the ability to set a context for a request.
func NewDeleteAutoprimaryParamsWithContext(ctx context.Context) *DeleteAutoprimaryParams {
	return &DeleteAutoprimaryParams{
		Context: ctx,
	}
}

// NewDeleteAutoprimaryParamsWithHTTPClient creates a new DeleteAutoprimaryParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteAutoprimaryParamsWithHTTPClient(client *http.Client) *DeleteAutoprimaryParams {
	return &DeleteAutoprimaryParams{
		HTTPClient: client,
	}
}

/* DeleteAutoprimaryParams contains all the parameters to send to the API endpoint
   for the delete autoprimary operation.

   Typically these are written to a http.Request.
*/
type DeleteAutoprimaryParams struct {

	/* IP.

	   IP address of autoprimary
	*/
	IP string

	/* Nameserver.

	   DNS name of the autoprimary
	*/
	Nameserver string

	/* ServerID.

	   The id of the server to delete the autoprimary from
	*/
	ServerID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the delete autoprimary params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteAutoprimaryParams) WithDefaults() *DeleteAutoprimaryParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete autoprimary params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteAutoprimaryParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the delete autoprimary params
func (o *DeleteAutoprimaryParams) WithTimeout(timeout time.Duration) *DeleteAutoprimaryParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete autoprimary params
func (o *DeleteAutoprimaryParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete autoprimary params
func (o *DeleteAutoprimaryParams) WithContext(ctx context.Context) *DeleteAutoprimaryParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete autoprimary params
func (o *DeleteAutoprimaryParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete autoprimary params
func (o *DeleteAutoprimaryParams) WithHTTPClient(client *http.Client) *DeleteAutoprimaryParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete autoprimary params
func (o *DeleteAutoprimaryParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithIP adds the ip to the delete autoprimary params
func (o *DeleteAutoprimaryParams) WithIP(ip string) *DeleteAutoprimaryParams {
	o.SetIP(ip)
	return o
}

// SetIP adds the ip to the delete autoprimary params
func (o *DeleteAutoprimaryParams) SetIP(ip string) {
	o.IP = ip
}

// WithNameserver adds the nameserver to the delete autoprimary params
func (o *DeleteAutoprimaryParams) WithNameserver(nameserver string) *DeleteAutoprimaryParams {
	o.SetNameserver(nameserver)
	return o
}

// SetNameserver adds the nameserver to the delete autoprimary params
func (o *DeleteAutoprimaryParams) SetNameserver(nameserver string) {
	o.Nameserver = nameserver
}

// WithServerID adds the serverID to the delete autoprimary params
func (o *DeleteAutoprimaryParams) WithServerID(serverID string) *DeleteAutoprimaryParams {
	o.SetServerID(serverID)
	return o
}

// SetServerID adds the serverId to the delete autoprimary params
func (o *DeleteAutoprimaryParams) SetServerID(serverID string) {
	o.ServerID = serverID
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteAutoprimaryParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param ip
	if err := r.SetPathParam("ip", o.IP); err != nil {
		return err
	}

	// path param nameserver
	if err := r.SetPathParam("nameserver", o.Nameserver); err != nil {
		return err
	}

	// path param server_id
	if err := r.SetPathParam("server_id", o.ServerID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
