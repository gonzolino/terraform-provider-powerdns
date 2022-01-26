package powerdns

import (
	"context"
	"errors"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/gonzolino/terraform-provider-powerdns/internal/powerdns/client"
	"github.com/gonzolino/terraform-provider-powerdns/internal/powerdns/client/zones"
	"github.com/gonzolino/terraform-provider-powerdns/internal/powerdns/models"
)

type Client struct {
	client   *apiclient.PowerDNSAuthoritativeHTTPAPI
	authInfo runtime.ClientAuthInfoWriter
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
	Name    string
	Type    string
	TTL     int64
	Records []string
}

func New(ctx context.Context, apiKey, serverHost, basePath, scheme string) *Client {
	transport := httptransport.New(serverHost, basePath, []string{scheme})

	return &Client{
		client:   apiclient.New(transport, strfmt.Default),
		authInfo: httptransport.APIKeyAuth("X-API-Key", "header", apiKey),
	}
}

func (pdns *Client) CreateZone(ctx context.Context, serverID string, zone *Zone) (*Zone, error) {
	if zone.Name == "" {
		return nil, errors.New("zone name is required")
	}
	if zone.Kind == "" {
		return nil, errors.New("zone kind is required")
	}
	zoneStruct := transformZoneToAPI(zone)

	responseRrsets := true
	params := zones.NewCreateZoneParamsWithContext(ctx).WithServerID(serverID).WithZoneStruct(zoneStruct).WithRrsets(&responseRrsets)

	resp, err := pdns.client.Zones.CreateZone(params, pdns.authInfo)
	if err != nil {
		return nil, err
	}

	return transformAPIToZone(resp.Payload), nil
}

func (pdns *Client) UpdateZone(ctx context.Context, serverID, zoneID string, zone *Zone) error {
	zoneStruct := transformZoneToAPI(zone)

	params := zones.NewPutZoneParamsWithContext(ctx).WithServerID(serverID).WithZoneID(zoneID).WithZoneStruct(zoneStruct)

	_, err := pdns.client.Zones.PutZone(params, pdns.authInfo)
	return err
}

func (pdns *Client) GetZone(ctx context.Context, serverID, zoneID string) (*Zone, error) {
	params := zones.NewListZoneParamsWithContext(ctx).WithServerID(serverID).WithZoneID(zoneID)

	resp, err := pdns.client.Zones.ListZone(params, pdns.authInfo)
	if err != nil {
		return nil, err
	}

	return transformAPIToZone(resp.Payload), nil
}

func (pdns *Client) DeleteZone(ctx context.Context, serverID, zoneID string) error {
	params := zones.NewDeleteZoneParamsWithContext(ctx).WithServerID(serverID).WithZoneID(zoneID)

	_, err := pdns.client.Zones.DeleteZone(params, pdns.authInfo)
	return err
}

func transformZoneToAPI(zone *Zone) *models.Zone {
	rrsets := make([]*models.RRSet, len(zone.RecordSets))
	for i, recordset := range zone.RecordSets {
		records := make([]*models.Record, len(recordset.Records))
		for j, record := range recordset.Records {
			records[j] = &models.Record{
				Content: &record,
			}
		}
		// TODO: we may have to copy record, else &recordset may all point to the same value
		rrsets[i] = &models.RRSet{
			Name:    &recordset.Name,
			Type:    &recordset.Type,
			TTL:     &recordset.TTL,
			Records: records,
		}
	}

	return &models.Zone{
		Name:    zone.Name,
		Kind:    zone.Kind,
		Dnssec:  zone.DNSSec,
		Serial:  zone.Serial,
		Masters: zone.Masters,
		Rrsets:  rrsets,
	}
}

func transformAPIToZone(zone *models.Zone) *Zone {
	recordsets := make([]RecordSet, len(zone.Rrsets))
	for i, rrset := range zone.Rrsets {
		records := make([]string, len(rrset.Records))
		for j, record := range rrset.Records {
			records[j] = *record.Content
		}
		recordsets[i] = RecordSet{
			Name:    *rrset.Name,
			Type:    *rrset.Type,
			TTL:     *rrset.TTL,
			Records: records,
		}
	}

	return &Zone{
		ID:         zone.ID,
		Name:       zone.Name,
		Kind:       zone.Kind,
		DNSSec:     zone.Dnssec,
		Serial:     zone.Serial,
		Masters:    zone.Masters,
		RecordSets: recordsets,
	}
}
