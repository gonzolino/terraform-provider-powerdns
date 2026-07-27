package powerdns

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	pdnsclient "github.com/gonzolino/terraform-provider-powerdns/internal/powerdns/client"
)

type Client struct {
	client pdnsclient.ClientWithResponsesInterface
}

// apiResponse is implemented by all generated *Response types and lets
// checkResponse inspect the HTTP status and any parsed error body without
// duplicating logic per endpoint.
type apiResponse interface {
	StatusCode() int
	GetJSONDefault() *pdnsclient.ErrorResponse
}

// checkResponse turns a non-transport-level API failure into a Go error.
// Unlike the previous go-openapi based client, the oapi-codegen client only
// returns a non-nil error for transport failures, so HTTP-level errors
// (4xx/5xx) must be detected explicitly via the response status code.
func checkResponse(resp apiResponse, wantStatus int) error {
	if resp.StatusCode() == wantStatus {
		return nil
	}
	if errResp := resp.GetJSONDefault(); errResp != nil {
		return fmt.Errorf("powerdns api error (status %d): %s", resp.StatusCode(), errResp.Error)
	}
	return fmt.Errorf("powerdns api error: unexpected status %d", resp.StatusCode())
}

type Zone struct {
	ID         string
	Name       string
	Kind       string
	DNSSec     bool
	Serial     int64
	Masters    []string
	RecordSets []RecordSet
}

type RecordSet struct {
	Name       string
	Type       string
	TTL        int64
	Changetype string
	Records    []string
}

func New(ctx context.Context, apiKey, serverHost, basePath, scheme string) (*Client, error) {
	baseURL := fmt.Sprintf("%s://%s%s", scheme, serverHost, basePath)

	authEditor := func(_ context.Context, req *http.Request) error {
		req.Header.Set("X-API-Key", apiKey)
		return nil
	}

	client, err := pdnsclient.NewClientWithResponses(baseURL, pdnsclient.WithRequestEditorFn(authEditor))
	if err != nil {
		return nil, fmt.Errorf("creating powerdns client: %w", err)
	}

	return &Client{client: client}, nil
}

func (pdns *Client) CreateZone(ctx context.Context, serverID string, zone *Zone) (*Zone, error) {
	if zone.Name == "" {
		return nil, errors.New("zone name is required")
	}
	if zone.Kind == "" {
		return nil, errors.New("zone kind is required")
	}

	responseRrsets := true
	params := &pdnsclient.CreateZoneParams{Rrsets: &responseRrsets}

	resp, err := pdns.client.CreateZoneWithResponse(ctx, serverID, params, transformZoneToAPI(zone))
	if err != nil {
		return nil, err
	}
	if err := checkResponse(resp, http.StatusCreated); err != nil {
		return nil, err
	}

	return transformAPIToZone(resp.JSON201), nil
}

func (pdns *Client) UpdateZone(ctx context.Context, serverID, zoneID string, zone *Zone) error {
	resp, err := pdns.client.PutZoneWithResponse(ctx, serverID, zoneID, transformZoneToAPI(zone))
	if err != nil {
		return err
	}

	return checkResponse(resp, http.StatusNoContent)
}

func (pdns *Client) GetZone(ctx context.Context, serverID, zoneID string) (*Zone, error) {
	resp, err := pdns.client.ListZoneWithResponse(ctx, serverID, zoneID, nil)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(resp, http.StatusOK); err != nil {
		return nil, err
	}

	return transformAPIToZone(resp.JSON200), nil
}

func (pdns *Client) DeleteZone(ctx context.Context, serverID, zoneID string) error {
	resp, err := pdns.client.DeleteZoneWithResponse(ctx, serverID, zoneID)
	if err != nil {
		return err
	}

	return checkResponse(resp, http.StatusNoContent)
}

func (pdns *Client) CreateRecordSet(ctx context.Context, serverID, zoneID string, recordSet *RecordSet) (*RecordSet, error) {
	rrset := transformRecordSetToAPI(recordSet)

	changeTypeReplace := "REPLACE"
	rrset.Changetype = &changeTypeReplace

	zone := pdnsclient.Zone{Rrsets: &[]pdnsclient.RRSet{rrset}}

	resp, err := pdns.client.PatchZoneWithResponse(ctx, serverID, zoneID, zone)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(resp, http.StatusNoContent); err != nil {
		return nil, err
	}

	return pdns.GetRecordSet(ctx, serverID, zoneID, recordSet.Name, recordSet.Type)
}

func (pdns *Client) UpdateRecordSet(ctx context.Context, serverID, zoneID string, recordSet *RecordSet) error {
	rrset := transformRecordSetToAPI(recordSet)

	changeTypeReplace := "REPLACE"
	rrset.Changetype = &changeTypeReplace

	zone := pdnsclient.Zone{Rrsets: &[]pdnsclient.RRSet{rrset}}

	resp, err := pdns.client.PatchZoneWithResponse(ctx, serverID, zoneID, zone)
	if err != nil {
		return err
	}

	return checkResponse(resp, http.StatusNoContent)
}

func (pdns *Client) GetRecordSet(ctx context.Context, serverID, zoneID, recordSetName, recordSetType string) (*RecordSet, error) {
	resp, err := pdns.client.ListZoneWithResponse(ctx, serverID, zoneID, nil)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(resp, http.StatusOK); err != nil {
		return nil, err
	}

	var allRrsets []pdnsclient.RRSet
	if resp.JSON200.Rrsets != nil {
		allRrsets = *resp.JSON200.Rrsets
	}

	rrs := []pdnsclient.RRSet{}
	for _, rrset := range allRrsets {
		if rrset.Name == recordSetName {
			rrs = append(rrs, rrset)
		}
	}

	switch len(rrs) {
	case 0:
		return nil, fmt.Errorf("record set '%s' not found", recordSetName)
	case 1:
		return transformAPIToRecordSet(&rrs[0]), nil
	default:
		if recordSetType == "" {
			return nil, errors.New("multiple record sets found with the same name, type required")
		}

		for _, rrset := range rrs {
			if rrset.Type == recordSetType {
				return transformAPIToRecordSet(&rrset), nil
			}
		}
		return nil, fmt.Errorf("record set '%s' with type '%s' not found", recordSetName, recordSetType)
	}
}

func (pdns *Client) DeleteRecordSet(ctx context.Context, serverID, zoneID string, recordSet *RecordSet) error {
	rrset := transformRecordSetToAPI(recordSet)

	changeTypeDelete := "DELETE"
	rrset.Changetype = &changeTypeDelete

	zone := pdnsclient.Zone{Rrsets: &[]pdnsclient.RRSet{rrset}}

	resp, err := pdns.client.PatchZoneWithResponse(ctx, serverID, zoneID, zone)
	if err != nil {
		return err
	}

	return checkResponse(resp, http.StatusNoContent)
}

func transformRecordSetToAPI(recordSet *RecordSet) pdnsclient.RRSet {
	records := make([]pdnsclient.Record, len(recordSet.Records))
	for i, record := range recordSet.Records {
		records[i] = pdnsclient.Record{
			Content: record,
		}
	}
	return pdnsclient.RRSet{
		Name:    recordSet.Name,
		Type:    recordSet.Type,
		Ttl:     int(recordSet.TTL),
		Records: records,
	}
}

func transformAPIToRecordSet(rrset *pdnsclient.RRSet) *RecordSet {
	records := make([]string, len(rrset.Records))
	for i, record := range rrset.Records {
		records[i] = record.Content
	}
	return &RecordSet{
		Name:    rrset.Name,
		Type:    rrset.Type,
		TTL:     int64(rrset.Ttl),
		Records: records,
	}
}

func transformZoneToAPI(zone *Zone) pdnsclient.Zone {
	rrsets := make([]pdnsclient.RRSet, len(zone.RecordSets))
	for i, recordset := range zone.RecordSets {
		records := make([]pdnsclient.Record, len(recordset.Records))
		for j, record := range recordset.Records {
			records[j] = pdnsclient.Record{
				Content: record,
			}
		}
		rrsets[i] = pdnsclient.RRSet{
			Name:    recordset.Name,
			Type:    recordset.Type,
			Ttl:     int(recordset.TTL),
			Records: records,
		}
	}

	kind := pdnsclient.ZoneKind(zone.Kind)
	serial := int(zone.Serial)

	return pdnsclient.Zone{
		Name:    &zone.Name,
		Kind:    &kind,
		Dnssec:  &zone.DNSSec,
		Serial:  &serial,
		Masters: &zone.Masters,
		Rrsets:  &rrsets,
	}
}

func transformAPIToZone(zone *pdnsclient.Zone) *Zone {
	var recordsets []RecordSet
	if zone.Rrsets != nil {
		recordsets = make([]RecordSet, len(*zone.Rrsets))
		for i, rrset := range *zone.Rrsets {
			records := make([]string, len(rrset.Records))
			for j, record := range rrset.Records {
				records[j] = record.Content
			}
			recordsets[i] = RecordSet{
				Name:    rrset.Name,
				Type:    rrset.Type,
				TTL:     int64(rrset.Ttl),
				Records: records,
			}
		}
	}

	result := &Zone{RecordSets: recordsets}
	if zone.Id != nil {
		result.ID = *zone.Id
	}
	if zone.Name != nil {
		result.Name = *zone.Name
	}
	if zone.Kind != nil {
		result.Kind = string(*zone.Kind)
	}
	if zone.Dnssec != nil {
		result.DNSSec = *zone.Dnssec
	}
	if zone.Serial != nil {
		result.Serial = int64(*zone.Serial)
	}
	if zone.Masters != nil {
		result.Masters = *zone.Masters
	}

	return result
}
