// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// TSIGKey TSIGKey
//
// A TSIG key that can be used to authenticate NOTIFY, AXFR, and DNSUPDATE queries.
//
// swagger:model TSIGKey
type TSIGKey struct {

	// The algorithm of the TSIG key
	Algorithm string `json:"algorithm,omitempty"`

	// The ID for this key, used in the TSIGkey URL endpoint.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// The Base64 encoded secret key, empty when listing keys. MAY be empty when POSTing to have the server generate the key material
	Key string `json:"key,omitempty"`

	// The name of the key
	Name string `json:"name,omitempty"`

	// Set to "TSIGKey"
	// Read Only: true
	Type string `json:"type,omitempty"`
}

// Validate validates this t s i g key
func (m *TSIGKey) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validate this t s i g key based on the context it is used
func (m *TSIGKey) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateType(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TSIGKey) contextValidateID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "id", "body", string(m.ID)); err != nil {
		return err
	}

	return nil
}

func (m *TSIGKey) contextValidateType(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "type", "body", string(m.Type)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TSIGKey) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TSIGKey) UnmarshalBinary(b []byte) error {
	var res TSIGKey
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
