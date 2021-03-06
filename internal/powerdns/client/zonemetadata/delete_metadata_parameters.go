// Code generated by go-swagger; DO NOT EDIT.

package zonemetadata

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

// NewDeleteMetadataParams creates a new DeleteMetadataParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteMetadataParams() *DeleteMetadataParams {
	return &DeleteMetadataParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteMetadataParamsWithTimeout creates a new DeleteMetadataParams object
// with the ability to set a timeout on a request.
func NewDeleteMetadataParamsWithTimeout(timeout time.Duration) *DeleteMetadataParams {
	return &DeleteMetadataParams{
		timeout: timeout,
	}
}

// NewDeleteMetadataParamsWithContext creates a new DeleteMetadataParams object
// with the ability to set a context for a request.
func NewDeleteMetadataParamsWithContext(ctx context.Context) *DeleteMetadataParams {
	return &DeleteMetadataParams{
		Context: ctx,
	}
}

// NewDeleteMetadataParamsWithHTTPClient creates a new DeleteMetadataParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteMetadataParamsWithHTTPClient(client *http.Client) *DeleteMetadataParams {
	return &DeleteMetadataParams{
		HTTPClient: client,
	}
}

/* DeleteMetadataParams contains all the parameters to send to the API endpoint
   for the delete metadata operation.

   Typically these are written to a http.Request.
*/
type DeleteMetadataParams struct {

	/* MetadataKind.

	   The kind of metadata
	*/
	MetadataKind string

	/* ServerID.

	   The id of the server to retrieve
	*/
	ServerID string

	/* ZoneID.

	   The id of the zone to retrieve
	*/
	ZoneID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the delete metadata params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteMetadataParams) WithDefaults() *DeleteMetadataParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete metadata params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteMetadataParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the delete metadata params
func (o *DeleteMetadataParams) WithTimeout(timeout time.Duration) *DeleteMetadataParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete metadata params
func (o *DeleteMetadataParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete metadata params
func (o *DeleteMetadataParams) WithContext(ctx context.Context) *DeleteMetadataParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete metadata params
func (o *DeleteMetadataParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete metadata params
func (o *DeleteMetadataParams) WithHTTPClient(client *http.Client) *DeleteMetadataParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete metadata params
func (o *DeleteMetadataParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithMetadataKind adds the metadataKind to the delete metadata params
func (o *DeleteMetadataParams) WithMetadataKind(metadataKind string) *DeleteMetadataParams {
	o.SetMetadataKind(metadataKind)
	return o
}

// SetMetadataKind adds the metadataKind to the delete metadata params
func (o *DeleteMetadataParams) SetMetadataKind(metadataKind string) {
	o.MetadataKind = metadataKind
}

// WithServerID adds the serverID to the delete metadata params
func (o *DeleteMetadataParams) WithServerID(serverID string) *DeleteMetadataParams {
	o.SetServerID(serverID)
	return o
}

// SetServerID adds the serverId to the delete metadata params
func (o *DeleteMetadataParams) SetServerID(serverID string) {
	o.ServerID = serverID
}

// WithZoneID adds the zoneID to the delete metadata params
func (o *DeleteMetadataParams) WithZoneID(zoneID string) *DeleteMetadataParams {
	o.SetZoneID(zoneID)
	return o
}

// SetZoneID adds the zoneId to the delete metadata params
func (o *DeleteMetadataParams) SetZoneID(zoneID string) {
	o.ZoneID = zoneID
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteMetadataParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param metadata_kind
	if err := r.SetPathParam("metadata_kind", o.MetadataKind); err != nil {
		return err
	}

	// path param server_id
	if err := r.SetPathParam("server_id", o.ServerID); err != nil {
		return err
	}

	// path param zone_id
	if err := r.SetPathParam("zone_id", o.ZoneID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
