package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPowerdnsRecordsetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccPowerdnsRecordsetDataSourceConfig("example.net.", "localhost", "example.net.", "NS"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerdns_recordset.test", "zone_id", "example.net."),
					resource.TestCheckResourceAttr("data.powerdns_recordset.test", "server_id", "localhost"),
					resource.TestCheckResourceAttr("data.powerdns_recordset.test", "name", "example.net."),
					resource.TestCheckResourceAttr("data.powerdns_recordset.test", "type", "NS"),
					resource.TestCheckResourceAttr("data.powerdns_recordset.test", "ttl", "1500"),
					resource.TestCheckResourceAttr("data.powerdns_recordset.test", "records.#", "2"),
				),
			},
			{
				Config:      testAccPowerdnsRecordsetDataSourceConfig("unknown.net.", "localhost", "example.net.", "NS"),
				ExpectError: regexp.MustCompile(`Unable to get record set 'example\.net\.' \(type 'NS'\): .*`),
			},
			{
				Config:      testAccPowerdnsRecordsetDataSourceConfig("example.net.", "unknownhost", "example.net.", "NS"),
				ExpectError: regexp.MustCompile(`Unable to get record set 'example\.net\.' \(type 'NS'\): .*`),
			},
			{
				Config:      testAccPowerdnsRecordsetDataSourceConfig("example.net.", "localhost", "unknown.example.net.", "A"),
				ExpectError: regexp.MustCompile(`Unable to get record set 'unknown\.example\.net\.' \(type 'A'\): .*`),
			},
			{
				Config:      testAccPowerdnsRecordsetDataSourceConfig("example.net.", "localhost", "example.net.", "unknown"),
				ExpectError: regexp.MustCompile(`Unable to get record set 'example\.net\.' \(type 'unknown'\): .*`),
			},
		},
	})
}

func testAccPowerdnsRecordsetDataSourceConfig(zoneId, serverId, name, typ string) string {
	return fmt.Sprintf(`
data "powerdns_recordset" "test" {
  zone_id = %[1]q
  server_id = %[2]q
  name = %[3]q
  type = %[4]q
}
`, zoneId, serverId, name, typ)
}
