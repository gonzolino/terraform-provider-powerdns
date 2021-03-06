// Code generated by go-swagger; DO NOT EDIT.

package search

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
	"github.com/go-openapi/swag"
)

// NewSearchDataParams creates a new SearchDataParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSearchDataParams() *SearchDataParams {
	return &SearchDataParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSearchDataParamsWithTimeout creates a new SearchDataParams object
// with the ability to set a timeout on a request.
func NewSearchDataParamsWithTimeout(timeout time.Duration) *SearchDataParams {
	return &SearchDataParams{
		timeout: timeout,
	}
}

// NewSearchDataParamsWithContext creates a new SearchDataParams object
// with the ability to set a context for a request.
func NewSearchDataParamsWithContext(ctx context.Context) *SearchDataParams {
	return &SearchDataParams{
		Context: ctx,
	}
}

// NewSearchDataParamsWithHTTPClient creates a new SearchDataParams object
// with the ability to set a custom HTTPClient for a request.
func NewSearchDataParamsWithHTTPClient(client *http.Client) *SearchDataParams {
	return &SearchDataParams{
		HTTPClient: client,
	}
}

/* SearchDataParams contains all the parameters to send to the API endpoint
   for the search data operation.

   Typically these are written to a http.Request.
*/
type SearchDataParams struct {

	/* Max.

	   Maximum number of entries to return
	*/
	Max int64

	/* ObjectType.

	   Type of data to search for, one of “all”, “zone”, “record”, “comment”
	*/
	ObjectType *string

	/* Q.

	   The string to search for
	*/
	Q string

	/* ServerID.

	   The id of the server to retrieve
	*/
	ServerID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the search data params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchDataParams) WithDefaults() *SearchDataParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the search data params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchDataParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the search data params
func (o *SearchDataParams) WithTimeout(timeout time.Duration) *SearchDataParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the search data params
func (o *SearchDataParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the search data params
func (o *SearchDataParams) WithContext(ctx context.Context) *SearchDataParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the search data params
func (o *SearchDataParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the search data params
func (o *SearchDataParams) WithHTTPClient(client *http.Client) *SearchDataParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the search data params
func (o *SearchDataParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithMax adds the max to the search data params
func (o *SearchDataParams) WithMax(max int64) *SearchDataParams {
	o.SetMax(max)
	return o
}

// SetMax adds the max to the search data params
func (o *SearchDataParams) SetMax(max int64) {
	o.Max = max
}

// WithObjectType adds the objectType to the search data params
func (o *SearchDataParams) WithObjectType(objectType *string) *SearchDataParams {
	o.SetObjectType(objectType)
	return o
}

// SetObjectType adds the objectType to the search data params
func (o *SearchDataParams) SetObjectType(objectType *string) {
	o.ObjectType = objectType
}

// WithQ adds the q to the search data params
func (o *SearchDataParams) WithQ(q string) *SearchDataParams {
	o.SetQ(q)
	return o
}

// SetQ adds the q to the search data params
func (o *SearchDataParams) SetQ(q string) {
	o.Q = q
}

// WithServerID adds the serverID to the search data params
func (o *SearchDataParams) WithServerID(serverID string) *SearchDataParams {
	o.SetServerID(serverID)
	return o
}

// SetServerID adds the serverId to the search data params
func (o *SearchDataParams) SetServerID(serverID string) {
	o.ServerID = serverID
}

// WriteToRequest writes these params to a swagger request
func (o *SearchDataParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param max
	qrMax := o.Max
	qMax := swag.FormatInt64(qrMax)
	if qMax != "" {

		if err := r.SetQueryParam("max", qMax); err != nil {
			return err
		}
	}

	if o.ObjectType != nil {

		// query param object_type
		var qrObjectType string

		if o.ObjectType != nil {
			qrObjectType = *o.ObjectType
		}
		qObjectType := qrObjectType
		if qObjectType != "" {

			if err := r.SetQueryParam("object_type", qObjectType); err != nil {
				return err
			}
		}
	}

	// query param q
	qrQ := o.Q
	qQ := qrQ
	if qQ != "" {

		if err := r.SetQueryParam("q", qQ); err != nil {
			return err
		}
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
