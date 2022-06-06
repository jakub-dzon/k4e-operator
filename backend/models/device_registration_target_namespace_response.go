// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DeviceRegistrationTargetNamespaceResponse device registration target namespace response
//
// swagger:model device-registration-target-namespace-response
type DeviceRegistrationTargetNamespaceResponse struct {

	// Exposes the error message generated at the backend when there is an error (example HTTP code 500).
	Message string `json:"message,omitempty"`

	// Contains namespace the device should be finally placed during registration.
	Namespace string `json:"namespace,omitempty"`
}

// Validate validates this device registration target namespace response
func (m *DeviceRegistrationTargetNamespaceResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this device registration target namespace response based on context it is used
func (m *DeviceRegistrationTargetNamespaceResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DeviceRegistrationTargetNamespaceResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DeviceRegistrationTargetNamespaceResponse) UnmarshalBinary(b []byte) error {
	var res DeviceRegistrationTargetNamespaceResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}