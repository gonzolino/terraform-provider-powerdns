// Code generated by go-swagger; DO NOT EDIT.

package zones

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// AxfrExportZoneReader is a Reader for the AxfrExportZone structure.
type AxfrExportZoneReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AxfrExportZoneReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewAxfrExportZoneOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewAxfrExportZoneOK creates a AxfrExportZoneOK with default headers values
func NewAxfrExportZoneOK() *AxfrExportZoneOK {
	return &AxfrExportZoneOK{}
}

/* AxfrExportZoneOK describes a response with status code 200, with default header values.

OK
*/
type AxfrExportZoneOK struct {
	Payload string
}

func (o *AxfrExportZoneOK) Error() string {
	return fmt.Sprintf("[GET /servers/{server_id}/zones/{zone_id}/export][%d] axfrExportZoneOK  %+v", 200, o.Payload)
}
func (o *AxfrExportZoneOK) GetPayload() string {
	return o.Payload
}

func (o *AxfrExportZoneOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
