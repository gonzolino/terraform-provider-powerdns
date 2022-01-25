package powerdns

import (
	"context"

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

func (pdns *Client) GetZone(ctx context.Context, serverID, zoneID string) (*Zone, error) {
	params := zones.NewListZoneParamsWithContext(ctx).WithServerID(serverID).WithZoneID(zoneID)

	resp, err := pdns.client.Zones.ListZone(params, pdns.authInfo)
	if err != nil {
		return nil, err
	}

	return transformAPIToZone(resp.Payload), nil
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
